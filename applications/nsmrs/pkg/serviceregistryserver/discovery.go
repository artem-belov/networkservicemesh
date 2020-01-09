// Copyright (c) 2020 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package serviceregistryserver

import (
	"context"

	"github.com/networkservicemesh/networkservicemesh/pkg/tools/spanhelper"

	"github.com/pkg/errors"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
)

type discoveryService struct {
	cache NSERegistryCache
}

func newDiscoveryService(cache NSERegistryCache) *discoveryService {
	return &discoveryService{
		cache: cache,
	}
}

func (d *discoveryService) FindNetworkService(ctx context.Context, request *registry.FindNetworkServiceRequest) (*registry.FindNetworkServiceResponse, error) {
	span := spanhelper.FromContext(ctx, "Nsmrs.FindNetworkService")
	defer span.Finish()
	logger := span.Logger()

	networkServiceEnpoints := d.cache.GetEndpoints(request.NetworkServiceName)
	if len(networkServiceEnpoints) == 0 {
		err := errors.Errorf("no NetworkService with name: %v", request.NetworkServiceName)
		logger.Errorf("Cannot find Network Service: %v", err)
		return nil, err
	}

	response := &registry.FindNetworkServiceResponse{
		NetworkService: &registry.NetworkService{
			Name:    request.NetworkServiceName,
			Payload: networkServiceEnpoints[0].NetworkService.Payload,
			Matches: networkServiceEnpoints[0].NetworkService.Matches,
		},
		NetworkServiceManagers: make(map[string]*registry.NetworkServiceManager),
		Payload:                networkServiceEnpoints[0].NetworkService.Payload,
	}

	for _, endpoint := range networkServiceEnpoints {
		response.NetworkServiceManagers[endpoint.NetworkServiceManager.Name] = endpoint.NetworkServiceManager
		response.NetworkServiceEndpoints = append(response.NetworkServiceEndpoints, endpoint.NetworkServiceEndpoint)
	}

	logger.Infof("FindNetworkService done: %v", response)

	return response, nil
}
