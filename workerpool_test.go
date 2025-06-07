package main

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"go.uber.org/goleak"
)

func TestUsualUsage(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	wPool, err := NewWorkerPool(ctx, 3, 6, 10)
	if err != nil {
		t.Fail()
	}

	buf1 := bytes.NewBufferString("aaa")
	buf2 := bytes.NewBufferString("bbb")
	buf3 := bytes.NewBufferString("ccc")
	buf4 := bytes.NewBufferString("ddd")
	buf5 := bytes.NewBufferString("eee")
	buf6 := bytes.NewBufferString("fff")
	buf7 := bytes.NewBufferString("ggg")
	buf8 := bytes.NewBufferString("kkk")
	buf9 := bytes.NewBufferString("eeee")

	wPool.Submit(buf1)
	wPool.Submit(buf2)
	wPool.Submit(buf3)

	wPool.AddWorker(ctx)

	wPool.Submit(buf4)
	wPool.Submit(buf5)

	wPool.DeleteWorker()
	wPool.DeleteWorker()
	wPool.DeleteWorker()
	wPool.DeleteWorker()
	t.Log("delete all workers")

	wPool.Submit(buf6)
	wPool.Submit(buf7)

	wPool.AddWorker(ctx)
	wPool.AddWorker(ctx)

	t.Log("add two workers")

	wPool.Submit(buf8)
	wPool.Submit(buf9)

	t.Log("start wait")
	wPool.Wait()
	t.Log("stop wait")

	t.Log("free pool")
	wPool.Free()
}

func TestUseAfterFree(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	wPool, err := NewWorkerPool(ctx, 2, 6, 10)
	if err != nil {
		t.Fail()
	}

	buf1 := bytes.NewBufferString("aaa")
	buf2 := bytes.NewBufferString("bbb")
	buf3 := bytes.NewBufferString("ccc")
	buf4 := bytes.NewBufferString("ddd")

	wPool.Submit(buf1)
	wPool.Submit(buf2)

	wPool.Free()

	wPool.AddWorker(ctx)
	wPool.AddWorker(ctx)

	wPool.Submit(buf3)
	wPool.Submit(buf4)

	t.Log("start wait")
	wPool.Wait()
	t.Log("stop wait")

	t.Log("free pool")
	wPool.Free()
}

func TestReachedMaxWorkers(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	wPool, err := NewWorkerPool(ctx, 2, 3, 10)
	if err != nil {
		t.Fail()
	}

	buf1 := bytes.NewBufferString("aaa")
	buf2 := bytes.NewBufferString("bbb")

	wPool.Submit(buf1)
	wPool.Submit(buf2)

	wPool.AddWorker(ctx)
	err = wPool.AddWorker(ctx)
	if !errors.Is(err, ErrSoMuchWorkers) {
		t.Fail()
	}

	t.Log("free pool")
	wPool.Free()
}

func TestCreateWPoolErrorBadWorkersNum(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	wPool, err := NewWorkerPool(ctx, 4, 3, 10)
	if !errors.Is(err, ErrBadWorkersNums) || wPool != nil {
		t.Fail()
	}
}

func TestCreateWPoolNullInitWorkers(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	wPool, err := NewWorkerPool(ctx, 0, 3, 10)
	if err != nil {
		t.Fail()
	}

	wPool.AddWorker(ctx)
	buf1 := bytes.NewBufferString("aaa")
	wPool.Submit(buf1)

	wPool.Free()
}

func TestDeleteNonExistsWorker(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	wPool, err := NewWorkerPool(ctx, 0, 3, 10)
	if err != nil {
		t.Fail()
	}

	err = wPool.DeleteWorker()
	if !errors.Is(err, ErrNoWorkers) {
		t.Fail()
	}
}
