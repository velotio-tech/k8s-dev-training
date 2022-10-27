package services

import (
	"context"
)

type Client interface {
	GetResource(name string) Resource
}

type Resource interface {
	Create(ctx context.Context) error
	List(ctx context.Context) error
	Update(ctx context.Context) error
	Delete(ctx context.Context) error
}
