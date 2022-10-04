package model

func (m *Default) Update() {
	if m.Err() != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Err(err)
		return
	}

	err = m.data.Hasher()
	if err != nil {
		m.Err(err)
		return
	}

	c := m.Store.Connect(m.User())
	c.Update(m.data)
	err = c.Err()
	if err != nil {
		m.Err(err)
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.Err(err)
		return
	}

	m.Hasher()
}
