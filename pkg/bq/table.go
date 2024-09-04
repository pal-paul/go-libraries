package bq

type Table struct {
	Name    string
	dataset *Dataset
}

func (b *Table) Append(data any) error {
	ins := b.dataset.client.Dataset(b.dataset.Name).Table(b.Name).Inserter()
	err := ins.Put(b.dataset.ctx, data)
	if err != nil {
		return err
	}
	return nil
}
