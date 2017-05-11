package uconfig_test

import "testing"

func TestDefaults(t *testing.T) {
	/*
		conf := struct {
			S string `default:"hello"`
			I int    `default:"2"`

			St struct {
				S string `default:"world"`
			}
		}{}

		err := uflag.Defaults(&conf)

		confs := uconf.Config(&conf).
			Visitor(defaults.Mapper).
			Walker(uconf.FileWalker("config.yaml", json.Unmarshal)).
			Visitor(uenv.Mapper, mySecretsMapper, uflag.Mapper)


		if err != nil {
			t.Fatal(err)
		}

		if conf.S != "hello" {
			t.Fatalf("Expected conf.s to be hello, got: %v", conf.S)
		}

		if conf.I != 2 {
			t.Fatalf("Expected conf.s to be hello, got: %v", conf.I)
		}

		if conf.St.S != "world" {
			t.Fatalf("Expected conf.s to be hello, got: %v", conf.St.S)
		}

		type SS struct {
			S string `default:"list"`
		}

		conf2 := struct {
			SS []SS
		}{SS: make([]SS, 2)}

		err = uflag.Defaults(&conf2)

		if err != nil {
			t.Fatal(err)
		}

		for _, ss := range conf2.SS {
			if ss.S != "list" {
				t.Fatalf("Expected conf2.SS.S to be `list`, got: %v", ss.S)
			}
		}

		log.Printf("conf:  %#v\nconf2: %#v", conf, conf2)
	*/
}
