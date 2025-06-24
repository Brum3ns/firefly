package option

import (
	"reflect"
)

// The structure validate is a alias for *Options in it's current state but holds all the validation/configuration functions.
// !IMPORTANT : The reciver functions *MUST* have the same name as the variable name (*if the variable is needed to be validated*). This makes it possible to loop over variables and easy validate them
// Note : (The reason why each option has its own function is for easier troubleshooting and also for easier development in the future)
type validate struct {
	opt            Option
	reflectValue   reflect.Value
	interfaceValue reflect.Value
	typ            reflect.Type
}

// validate the user input to be correct before starting any future processes
func Validate(opt Option) (Option, error) {
	valid := &validate{opt: opt}
	valid.reflectValue = reflect.ValueOf(opt)

	if valid.reflectValue.Kind() == reflect.Ptr {
		valid.interfaceValue = valid.reflectValue.Elem()
	} else {
		valid.interfaceValue = valid.reflectValue
	}

	valid.typ = valid.interfaceValue.Type()

	for i := 0; i < valid.interfaceValue.NumField(); i++ {
		fValue := valid.interfaceValue.FieldByName(valid.typ.Field(i).Name)
		fTyp := fValue.Type()

		// Extract the option group variables, then extract all (options/flags) within the group:
		if fTyp.Kind() == reflect.Struct {
			for j := 0; j < fValue.NumField(); j++ {
				item := fTyp.Field(j)

				// Validation error detected for user input, return error to the user screen:
				if exist, ok := valid.MethodCall(item.Name); exist && !ok {
					return opt, getErrorByCode(item.Tag.Get("errcode"))
				}
			}
		}
	}

	return opt, nil
}

// Declare static values to be global:
// !Note : (This variable should NEVER be changed once defined)
/* func (conf *validate) setGlobal() bool {
	global.RANDOM_INSERT = conf.opt.Random
	global.VERBOSE = conf.opt.Verbose
	global.PAYLOAD_PATTERN = conf.opt.PayloadPattern
	return true
} */

// Call a method within the validate struct.
// Return if the method exist and if the method execution status, otherwise if the method wasen't found return double false values
func (conf *validate) MethodCall(name string) (bool, bool) {
	if v := reflect.ValueOf(conf).MethodByName(name); v.IsValid() {
		return true, v.Call(nil)[0].Interface() == true
	}
	return false, false
}

/*
func (valid *validate) Output() bool {
	return false
} */

//Configure each option in it's own method for easy code read and development in the feature
/*
func (conf *validate) ReqRaw() bool {
	return true
}

func (conf *validate) Encode() bool {
	return true
}
func (conf *validate) Tamper() bool {
	return true
}

func (conf *validate) Technique() bool {
	return true
}

func (conf *validate) Output() bool {
	return !files.FileExist(conf.opt.Output) || conf.opt.Overwrite
}

func (conf *validate) MaxIdleConns() bool {
	return conf.opt.MaxIdleConns > 0
}
func (conf *validate) MaxIdleConnsPerHost() bool {
	return conf.opt.MaxIdleConnsPerHost > 0
}
func (conf *validate) MaxConnsPerHost() bool {
	return conf.opt.MaxConnsPerHost > 0
}

func (conf *validate) AutoDetectParams() bool { //Can be ignored ATM
	return len(conf.opt.AutoParamRules) == 0 || conf.opt.AutoParamRules == "append" || conf.opt.AutoParamRules == "replace"
}

func (conf *validate) URLs() bool {
	return len(conf.opt.URLs) > 0
}

func (conf *validate) PayloadReplace() bool {
	return len(conf.opt.PayloadReplace) == 0 || (len(conf.opt.PayloadReplace) > 0 && strings.Contains(conf.opt.PayloadReplace, " => "))
}

func (conf *validate) MatchMode() bool {
	mode := strings.ToLower(conf.opt.MatchMode)
	return mode == "or" || mode == "and"
}
func (conf *validate) FilterMode() bool {
	mode := strings.ToLower(conf.opt.FilterMode)
	return mode == "or" || mode == "and"
}
func (conf *validate) Scheme() bool {
	return len(conf.opt.Scheme) > 0
}
func (conf *validate) Methods() bool {
	return len(conf.opt.Methods) > 0
}
func (conf *validate) Scanner() bool {
	return conf.opt.ThreadsScanner > 0
}
func (conf *validate) Threads() bool {
	return conf.opt.Threads > 0
}
func (conf *validate) Insert() bool {
	return len(conf.opt.InsertKeyword) > 0
}
func (conf *validate) Delay() bool {
	return conf.opt.Delay >= 0
}
func (conf *validate) Timeout() bool {
	return conf.opt.Timeout >= 0
}

func (conf *validate) ThreadsExtract() bool {
	return conf.opt.ThreadsExtract >= 0
}
func (conf *validate) VerifyAmount() bool {
	return conf.opt.VerifyAmount > 0
}
func (conf *validate) WordlistPaths() bool {
	return len(conf.opt.wordlistPath) > 0 && len(conf.opt.WordlistPaths) > 0
}
*/
