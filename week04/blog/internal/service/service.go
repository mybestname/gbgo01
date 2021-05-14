package service

import (
	pb "blog/api/blog/v1"
	"blog/internal/biz"

	"github.com/google/wire"
	"my/log"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewBlogService)

type BlogService struct {
	pb.UnimplementedBlogServiceServer

	log *log.Helper

	article *biz.ArticleUsecase
}
