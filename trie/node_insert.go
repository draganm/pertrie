package trie

import "fmt"

func (n Node) Put(key, value []byte) error {
	prefix, err := n.Prefix()
	if err != nil {
		return fmt.Errorf("could not get prefix: %w", err)
	}

	if prefix.IsEmpty() {
		nv := n.Value()
		switch nv.Which() {
		case Node_value_Which_nil:
			// this is good
			prefix, err := n.NewPrefix()
			if err != nil {
				return fmt.Errorf("could not create prefix: %w", err)
			}
			err = prefix.SetData(key)
			if err != nil {
				return fmt.Errorf("could not set prefix data: %w", err)
			}
			content, err := nv.NewContent()
			if err != nil {
				return fmt.Errorf("could not create value content: %w", err)
			}
			contentVal, err := content.Value().NewValue()
			if err != nil {
				return fmt.Errorf("could not set contentValue: %w", err)
			}
			// TODO: support for large(r) values
			contentVal.SetSize(uint64(len(value)))
			vs, err := contentVal.NewFirstSegment()
			if err != nil {
				return fmt.Errorf("could not create first segment: %w", err)
			}
			err = vs.SetData(value)
			if err != nil {
				return fmt.Errorf("could not set first segment value: %w", err)
			}

			return nil

		case Node_value_Which_blockRef:
			return fmt.Errorf("inserting when blockref is set is not yet supported")
		case Node_value_Which_content:
			return fmt.Errorf("inserting when content is set is not yet supported")
		default:
			return fmt.Errorf("this case should never be reached")
		}
	}

	return nil

}
