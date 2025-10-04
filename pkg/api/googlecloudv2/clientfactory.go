// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package googlecloudv2

import (
	"context"

	"google.golang.org/api/composer/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/gkehub/v1"
	"google.golang.org/api/logging/v2"
	"google.golang.org/api/option"
)

// ClientFactoryContextModifiers defines a function type for modifying the context
// before creating a Google Cloud client.
type ClientFactoryContextModifiers = func(ctx context.Context, container ResourceContainer) (context.Context, error)

// ClientFactoryOptionsModifiers defines a function type for modifying the client options
// before creating a Google Cloud client.
type ClientFactoryOptionsModifiers = func(ctx []option.ClientOption, container ResourceContainer) ([]option.ClientOption, error)

// ClientFactoryOption defines a function type for configuring a ClientFactory.
type ClientFactoryOption = func(s *clientFactory) error

// clientFactory generates a context used for generating the google cloud client.
type clientFactory struct {
	clientOptions    []ClientFactoryOptionsModifiers
	contextModifiers []ClientFactoryContextModifiers
}

// NewClientFactory creates a new ClientFactory with the given options.
// It applies each option to the factory and returns an error if any option fails.
func NewClientFactory(options ...ClientFactoryOption) (*clientFactory, error) {
	var factory = &clientFactory{}
	for _, opt := range options {
		err := opt(factory)
		if err != nil {
			return nil, err
		}
	}
	return factory, nil
}

// context applies all registered context modifiers to the given context for the given resource container.
func (s *clientFactory) context(ctx context.Context, container ResourceContainer) (context.Context, error) {
	for _, modifier := range s.contextModifiers {
		var err error
		ctx, err = modifier(ctx, container)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

// options applies all registered client options modifiers to an initial set of client options
// for the given resource container. It returns the modified client options or an error if any modifier fails.
func (s *clientFactory) options(container ResourceContainer) ([]option.ClientOption, error) {
	var options []option.ClientOption
	for _, modifier := range s.clientOptions {
		var err error
		options, err = modifier(options, container)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

// ContainerService returns the client for container.googleapis.com from given context and the resource container.
func (s *clientFactory) ContainerService(ctx context.Context, c ResourceContainer) (*container.Service, error) {
	ctx, err := s.context(ctx, c)
	if err != nil {
		return nil, err
	}
	options, err := s.options(c)
	if err != nil {
		return nil, err
	}

	return container.NewService(ctx, options...)
}

// GKEOnPremService returns the client for gkeonprem.googleapis.com from given context and the resource container.
func (s *clientFactory) GKEOnPremService(ctx context.Context, c ResourceContainer) (*gkehub.Service, error) {
	ctx, err := s.context(ctx, c)
	if err != nil {
		return nil, err
	}
	options, err := s.options(c)
	if err != nil {
		return nil, err
	}

	return gkehub.NewService(ctx, options...)
}

// ComposerService returns the client for composer.googleapis.com from given context and the resource container.
func (s *clientFactory) ComposerService(ctx context.Context, c ResourceContainer) (*composer.Service, error) {
	ctx, err := s.context(ctx, c)
	if err != nil {
		return nil, err
	}
	options, err := s.options(c)
	if err != nil {
		return nil, err
	}

	return composer.NewService(ctx, options...)
}

// LoggingService returns the client for logging.googleapis.com from given context and the resource container.
func (s *clientFactory) LoggingService(ctx context.Context, c ResourceContainer) (*logging.Service, error) {
	ctx, err := s.context(ctx, c)
	if err != nil {
		return nil, err
	}
	options, err := s.options(c)
	if err != nil {
		return nil, err
	}

	return logging.NewService(ctx, options...)
}
