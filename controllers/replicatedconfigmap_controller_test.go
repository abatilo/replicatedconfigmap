package controllers

import (
	"context"
	"time"

	rcmv1beta1 "github.com/abatilo/replicatedconfigmap/api/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("ReplicatedConfigMap Controller", func() {

	const timeout = time.Second * 10
	const interval = time.Second * 1

	Context("Create ConfigMaps", func() {
		It("Should create successfully", func() {

			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "namespace",
					Annotations: map[string]string{
						"rcm-sync": "true",
					},
				},
			}

			rcm := &rcmv1beta1.ReplicatedConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "replicatedconfigmap",
				},
				Spec: rcmv1beta1.ReplicatedConfigMapSpec{
					Name:     "configmap",
					Metadata: "rcm-metadata",
					Data: map[string]string{
						"foo": "bar",
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), namespace)).Should(Succeed())
			Expect(k8sClient.Create(context.Background(), rcm)).Should(Succeed())

			By("Expecting configmap to be created")
			Eventually(func() string {
				cm := &corev1.ConfigMap{}
				_ = k8sClient.Get(context.Background(), types.NamespacedName{
					Name:      "configmap",
					Namespace: "namespace",
				}, cm)
				return cm.Data["foo"]
			}, timeout, interval).Should(Equal("bar"))

			rcm.Spec.Data["foo"] = "baz"
			Expect(k8sClient.Update(context.Background(), rcm)).Should(Succeed())

			By("Expecting configmap to be updated")
			Eventually(func() string {
				cm := &corev1.ConfigMap{}
				_ = k8sClient.Get(context.Background(), types.NamespacedName{
					Name:      "configmap",
					Namespace: "namespace",
				}, cm)
				return cm.Data["foo"]
			}, timeout, interval).ShouldNot(Equal("bar"))
			Eventually(func() string {
				cm := &corev1.ConfigMap{}
				_ = k8sClient.Get(context.Background(), types.NamespacedName{
					Name:      "configmap",
					Namespace: "namespace",
				}, cm)
				return cm.Data["foo"]
			}, timeout, interval).Should(Equal("baz"))
		})
	})
})
