package provider

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// This function is from  dynamic fake
// and is licensed under the Apache License, Version 2.0 (the "License");
// https://github.com/kubernetes/kubernetes/blob/bf9a868e8ea3d3a8fa53cbb22f566771b3f8068b/staging/src/k8s.io/client-go/dynamic/fake/simple_test.go#L30
func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
			"spec": name,
		},
	}
}
