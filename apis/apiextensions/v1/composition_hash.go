package v1

import (
	"crypto/sha256"
	"fmt"

	"sigs.k8s.io/yaml"
)

func (c Composition) Hash() string {
	h := sha256.New()

	y, err := yaml.Marshal(c.ObjectMeta.Labels)
	if err != nil {
		return "labels failure during hashing"
	}
	a, err := yaml.Marshal(c.ObjectMeta.Annotations)
	if err != nil {
		return "annotations failure during hashing"
	}

	s, err := yaml.Marshal(c.Spec)
	if err != nil {
		return "spec failure during hashing"
	}

	y = append(y, a...)
	y = append(y, s...)
	_, _ = h.Write(y)

	return fmt.Sprintf("%x", h.Sum(nil))
}
