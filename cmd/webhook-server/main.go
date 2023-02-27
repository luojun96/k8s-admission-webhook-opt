package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	podResource = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

// applyAdmissionMutationOpertation implements the logic of admission controller webhook. For every pod that
// is created (outside of Kubernetes namespaces), it first checks if `runAsNonRoot` is set. If it is not, it
// is set to a default value of `false`. Furthermore, if `runAsUser` is not set (and `runAsNonRoot` was not
// initially set), it defaults `runAsUser` to a value of 1234.
//
// To demonstrate how requests can be rejected, this webhook further validates that the `runAsNonRoot`
// setting does not conflict with the `runAsUser` setting - i.e., if the former is set to `true`, the
// latter must not be `0`. Note that we combine both the setting of defaults and the check for potential
// conflicts in one webhook; ideally, the latter would be performed in a validating webhook admission controller.
func applyAdmissionMutationOpertation(req *admission.AdmissionRequest) ([]patchOperation, error) {
	if req.Resource != podResource {
		log.Printf("expect resource to be %s", podResource)
		return nil, nil
	}

	raw := req.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	// Retrieve the `runAsNonRoot` and `runAsUser` values.
	var runAsNonRoot *bool
	var runAsUser *int64
	if pod.Spec.SecurityContext != nil {
		runAsNonRoot = pod.Spec.SecurityContext.RunAsNonRoot
		runAsUser = pod.Spec.SecurityContext.RunAsUser
	}

	var patches []patchOperation
	if runAsNonRoot == nil {
		patches = append(patches, patchOperation{
			Op:    "add",
			Path:  "/spec/securityContext/runAsNonRoot",
			Value: runAsUser == nil || *runAsUser != 0,
		})

		if runAsUser == nil {
			patches = append(patches, patchOperation{
				Op:    "add",
				Path:  "/spec/securityContext/runAsUser",
				Value: 1234,
			})
		}
	} else if *runAsNonRoot == true && (runAsUser != nil && *runAsUser == 0) {
		return nil, errors.New("runAsNonRoot specified, but runAsUser set to 0 (the root user)")
	}

	return patches, nil
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applyAdmissionMutationOpertation))
	mux.Handle("/ready", readinessCheckHandler)
	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
