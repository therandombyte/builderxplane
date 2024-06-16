package composition

import (
	"fmt"

	v1 "github.com/therandombyte/builderxplane/apis/apiextensions/v1"
	meta "github.com/therandombyte/builderxplane/pkg/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewCompositionRevision(c *v1.Composition, revision int64) *v1.CompositionRevision {
	hash := c.Hash()
	if len(hash) >= 63 {
		hash = hash[0:63]
	}
	nameSuffix := hash
	if len(nameSuffix) >= 7 {
		nameSuffix = nameSuffix[0:7]
	}

	cr := &v1.CompositionRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", c.GetName(), nameSuffix),
			Labels: map[string]string{
				v1.LabelCompositionName: c.GetName(),
				v1.LabelCompositionHash: hash,
			},
		},
		Spec: NewCompostionRevisionSpec(c.Spec, revision),
	}
	ref := meta.TypedReferenceTo(c, v1.CompositionGroupVersionKind)
	meta.AddOwnerReference(cr, meta.AsController(ref))

	for k, v := range c.GetLabels() {
		cr.ObjectMeta.Labels[k] = v
	}

	return cr
}

func NewCompostionRevisionSpec(cs v1.CompositionSpec, revision int64) v1.CompositionRevisionSpec {
	conv := v1.GeneratedRevisionSpecConverter{}
	rs := conv.ToRevisionSpec(cs)
	rs.Revision = revision
	return rs
}
