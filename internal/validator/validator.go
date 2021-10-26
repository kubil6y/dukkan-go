package validator

type Validator struct {
	Errors map[string]string
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(condition bool, key, message string) {
	if !condition {
		v.AddError(key, message)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// NOTE helper functions: https://github.com/asaskevich/govalidator
