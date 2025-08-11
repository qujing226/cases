package ErrGroup

import "golang.org/x/sync/errgroup"

// errGroup 用来并发做两件不会互相干扰的事，但是需要都能完成
func doBoth() error {
	var (
		eg errgroup.Group
		a  A
		b  B
	)
	eg.Go(func() error {
		var er error
		a, er = a.Do()
		return er
	})

	eg.Go(func() error {
		var er error
		b, er = b.Do()
		return er
	})
	err := eg.Wait()
	return err
}

// errGroup 为一个对象并发查询字段
func findForA(a *A) error {
	var er errgroup.Group
	er.Go(func() error {
		var er error
		a.name, er = a.FindName()
		return er
	})
	er.Go(func() error {
		var er error
		a.age, er = a.FindAge()
		return er
	})
	return er.Wait()
}

type A struct {
	name string
	age  int
}

func (a *A) Do() (A, error) {
	return A{}, nil
}

func (a *A) FindName() (string, error) {
	// DAO \ cache
	return "A", nil
}

func (a *A) FindAge() (int, error) {
	// DAO \ cache
	return 18, nil
}

type B struct {
}

func (a *B) Do() (B, error) {
	return B{}, nil
}
