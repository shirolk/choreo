// Copyright 2025 The OpenChoreo Authors
// SPDX-License-Identifier: Apache-2.0

package choreoctl

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	openchoreov1alpha1 "github.com/openchoreo/openchoreo/api/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
	logger logr.Logger
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(openchoreov1alpha1.AddToScheme(scheme))
	setupLogger()
}

func setupLogger() {
	opts := zap.Options{
		Development: false,
	}
	logger = zap.New(zap.UseFlagOptions(&opts))
	log.SetLogger(logger)
}
