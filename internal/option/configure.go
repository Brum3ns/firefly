package option

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/pkg/files"
)

// The structure configure is a alias for *Options in it's current state but holds all the validation/configuration functions.
// !IMPORTANT : The reciver functions *MUST* have the same name as the variable name (*if the variable is needed to be validated*). This makes it possible to loop over variables and easy configure them
// Note : (The reason why each option has its own function is for easier troubleshooting and also for easier development in the future)
type configure struct {
	opt            *Options
	reflectValue   reflect.Value
	interfaceValue reflect.Value
	typ            reflect.Type
}

// validate the user input to be correct before starting any future processes
func Configure(opt *Options) (*Options, int) {
	conf := &configure{opt: opt}
	conf.reflectValue = reflect.ValueOf(conf.opt)
	conf.interfaceValue = conf.reflectValue.Elem()
	conf.typ = conf.interfaceValue.Type()

	for i := 0; i < conf.interfaceValue.NumField(); i++ {
		fValue := conf.interfaceValue.FieldByName(conf.typ.Field(i).Name)
		fTyp := fValue.Type()

		//Extract the option group variables, then extract all (options/flags) within the group:
		if fTyp.Kind() == reflect.Struct {
			for i := 0; i < fValue.NumField(); i++ {
				item := fTyp.Field(i)

				//Validation error detected for user input, return error to the user screen:
				if exist, ok := conf.MethodCall(item.Name); exist && !ok {
					if errcode, ok := strconv.Atoi(item.Tag.Get("errorcode")); ok == nil {
						return nil, errcode
					} else {
						log.Panicf("can't convert errorcode value \"%v\" for flag \"%s\".\n", errcode, item.Name)
					}
				}
			}
		}
	}

	//Define global variables (only none sensitive)
	conf.setGlobal()

	return conf.opt, 0
}

// Declare static values to be global:
// !Note : (This variable should NEVER be changed once defined)
func (conf *configure) setGlobal() bool {
	global.RANDOM_INSERT = conf.opt.Random
	global.VERBOSE = conf.opt.Verbose
	global.PAYLOAD_PATTERN = conf.opt.PayloadPattern
	return true
}

// Call a method within the validate struct.
// Return if the method exist and if the method execution status, otherwise if the method wasen't found return double false values
func (conf *configure) MethodCall(name string) (bool, bool) {
	if v := reflect.ValueOf(conf).MethodByName(name); v.IsValid() {
		return true, v.Call(nil)[0].Interface() == true
	}
	return false, false
}

// //////////Configure Options////////// //
//Configure each option in it's own method for easy code read and development in the feature

func (conf *configure) ReqRaw() bool {
	return true
}

func (conf *configure) Encode() bool {
	return true
}
func (conf *configure) Tamper() bool {
	return true
}

func (conf *configure) Technique() bool {
	return true
}

func (conf *configure) Output() bool {
	return !files.FileExist(conf.opt.Output) || conf.opt.Overwrite
}

func (conf *configure) MaxIdleConns() bool {
	return conf.opt.MaxIdleConns > 0
}
func (conf *configure) MaxIdleConnsPerHost() bool {
	return conf.opt.MaxIdleConnsPerHost > 0
}
func (conf *configure) MaxConnsPerHost() bool {
	return conf.opt.MaxConnsPerHost > 0
}

/* func (conf *configure) AutoDetectParams() bool { //Can be ignored ATM
	return len(conf.opt.AutoParamRules) == 0 || conf.opt.AutoParamRules == "append" || conf.opt.AutoParamRules == "replace"
} */

func (conf *configure) PayloadReplace() bool {
	return len(conf.opt.PayloadReplace) == 0 || (len(conf.opt.PayloadReplace) > 0 && strings.Contains(conf.opt.PayloadReplace, " => "))
}
func (conf *configure) URLs() bool {
	return len(conf.opt.URLs) > 0
}
func (conf *configure) MatchMode() bool {
	mode := strings.ToLower(conf.opt.MatchMode)
	return mode == "or" || mode == "and"
}
func (conf *configure) FilterMode() bool {
	mode := strings.ToLower(conf.opt.FilterMode)
	return mode == "or" || mode == "and"
}
func (conf *configure) Scheme() bool {
	return len(conf.opt.Scheme) > 0
}
func (conf *configure) Methods() bool {
	return len(conf.opt.Methods) > 0
}
func (conf *configure) Scanner() bool {
	return conf.opt.ThreadsScanner > 0
}
func (conf *configure) Threads() bool {
	return conf.opt.Threads > 0
}
func (conf *configure) Insert() bool {
	return len(conf.opt.InsertKeyword) > 0
}
func (conf *configure) Delay() bool {
	return conf.opt.Delay >= 0
}
func (conf *configure) Timeout() bool {
	return conf.opt.Timeout >= 0
}

func (conf *configure) ThreadsExtract() bool {
	return conf.opt.ThreadsExtract >= 0
}
func (conf *configure) VerifyAmount() bool {
	return conf.opt.VerifyAmount > 0
}
func (conf *configure) WordlistPaths() bool {
	return len(conf.opt.wordlistPath) > 0 && len(conf.opt.WordlistPaths) > 0
}
