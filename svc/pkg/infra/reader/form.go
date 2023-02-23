package reader

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"ynufes-mypage-backend/svc/pkg/domain/model/event"
	"ynufes-mypage-backend/svc/pkg/domain/model/form"
	entity "ynufes-mypage-backend/svc/pkg/infra/entity/form"
)

type Form struct {
	collection *firestore.CollectionRef
}

func NewForm(client *firestore.Client) *Form {
	return &Form{
		collection: client.Collection(entity.FormCollectionName),
	}
}

func (f Form) GetByID(ctx context.Context, id form.ID) (*form.Form, error) {
	snap, err := f.collection.Doc(id.ExportID()).Get(ctx)
	if err != nil {
		return nil, err
	}
	var e entity.Form
	err = snap.DataTo(&e)
	if err != nil {
		return nil, err
	}
	e.ID = snap.Ref.ID
	m, err := e.ToModel()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (f Form) ListByEventID(ctx context.Context, eventID event.ID) ([]form.Form, error) {
	iter := f.collection.Where("event_id", "==", eventID.GetValue()).Documents(ctx)
	var forms []form.Form
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}
		var e entity.Form
		err = doc.DataTo(&e)
		if err != nil {
			return nil, fmt.Errorf("failed to decode snap into entity.Form in ListByEventID: %w", err)
		}
		e.ID = doc.Ref.ID
		m, err := e.ToModel()
		if err != nil {
			return nil, fmt.Errorf("failed to convert entity to model in ListByEventID: %w", err)
		}
		forms = append(forms, *m)
	}
	return forms, nil
}
