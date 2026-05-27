package utils

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/sync/errgroup"
)

var (
	CST = time.FixedZone("CST", 8*60*60)
)

func GetUuid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func TernaryOperator[T any](expr bool, trueVal, falseVal T) T {
	if expr {
		return trueVal
	} else {
		return falseVal
	}
}

func ContextCopy(ctx context.Context) context.Context {
	return zipkin.NewContext(context.Background(), zipkin.SpanOrNoopFromContext(ctx))
}

func GenerateRandomIntList(size int, edge1, edge2 int) []int {
	list := make([]int, size, size)
	if edge1 > edge2 {
		edge1, edge2 = edge2, edge1
	}
	for i := 0; i < size; i++ {
		if edge1 == edge2 {
			list[i] = edge1
		} else {
			list[i] = rand.Intn(edge2-edge1) + edge1
		}
	}
	return list
}

func Max[T NUMBER](m, n T) T {
	return TernaryOperator(m > n, m, n)
}

func Min[T NUMBER](m, n T) T {
	return TernaryOperator(m < n, m, n)
}

func CurrentLimit(limit int, workList []func() error, enableError bool) error {
	if limit <= 0 || len(workList) == 0 {
		return nil
	}

	eg := new(errgroup.Group)
	limitCh := make(chan struct{}, limit)

	for _, w := range workList {
		work := w
		eg.Go(func() error {
			limitCh <- struct{}{}
			defer func() { <-limitCh }()
			if err := work(); err != nil {
				if !enableError {
					return err
				}
			}
			return nil
		})
	}
	return eg.Wait()
}

func Paginated(total int64, pageSize int, work func(int) error) error {
	for i := 0; i <= int(total/int64(pageSize)); i++ {
		if err := work(i); err != nil {
			return err
		}
	}
	return nil
}
