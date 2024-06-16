// to deal with k8s object metadata
package meta

import (
	xpv1 "github.com/therandombyte/builderxplane/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TypedReferenceTo(o metav1.Object, of schema.GroupVersionKind) *xpv1.TypedReference {
	v, k := of.ToAPIVersionAndKind()
	return &xpv1.TypedReference{
		APIVersion: v,
		Kind:       k,
		Name:       o.GetName(),
		UID:        o.GetUID(),
	}
}

func AddOwnerReference(o metav1.Object, r metav1.OwnerReference) {
	refs := o.GetOwnerReferences()
	for i := range refs {
		if refs[i].UID == r.UID {
			refs[i] = r
			o.SetOwnerReferences(refs)
			return
		}
	}
	o.SetOwnerReferences(append(refs, r))
}

// AsController converts the supplied object reference to a controller
// reference. You may also consider using metav1.NewControllerRef.
func AsController(r *xpv1.TypedReference) metav1.OwnerReference {
	t := true
	ref := AsOwner(r)
	ref.Controller = &t
	ref.BlockOwnerDeletion = &t
	return ref
}

// AsOwner converts the supplied object reference to an owner reference.
func AsOwner(r *xpv1.TypedReference) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Name:       r.Name,
		UID:        r.UID,
	}
}
