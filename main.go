package main

import (
	"core/client"
	"core/copy_db"
	"core/server"
	"core/transform"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/valyala/fasthttp"
)

type CitiItem struct {
	id      int64  `reindex:"id,hash,pk"`
	country int64  `reindex:"country,hash"`
	name    string `reindex:"name,hash"`
}

type Config struct {
	Server struct {
		Name    string
		Address string
	}
	DB struct {
		Host string
		Port string
	}
	VMName string
}

var Conf Config

func PreparingFiles() {

}
func gui_pro() {
	gtk.Init(nil)

	var setTemp int
	var path string

	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	err = b.AddFromFile("client_copy.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	labelAddressSetObj, err := b.GetObject("label_address_set")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	labelAddressSet := labelAddressSetObj.(*gtk.Label)

	labelNameSetObj, err := b.GetObject("label_name_set")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	labelNameSet := labelNameSetObj.(*gtk.Label)

	btnChoseDbObj, err := b.GetObject("btn_chose_db")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnChoseDb := btnChoseDbObj.(*gtk.Button)
	btnChoseDb.Connect("clicked", func() {
		setTemp = 1
		labelNameSet.SetText("Хост")
		labelAddressSet.SetText("Порт")
	})

	btnChoseVMObj, err := b.GetObject("btn_chose_vm")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnChoseVM := btnChoseVMObj.(*gtk.Button)
	btnChoseVM.Connect("clicked", func() {
		setTemp = 2
		labelNameSet.SetText("Имя машины")
		labelAddressSet.SetText("")
	})

	btnChoseConnectionObj, err := b.GetObject("btn_chose_connection")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnChoseConnection := btnChoseConnectionObj.(*gtk.Button)
	btnChoseConnection.Connect("clicked", func() {
		setTemp = 0
		labelNameSet.SetText("Название")
		labelAddressSet.SetText("Адрес")
	})
	//главное окно
	winMainObj, err := b.GetObject("window_main")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	winMain := winMainObj.(*gtk.Window)
	winMain.Connect("destroy", func() { gtk.MainQuit() })
	winMain.ShowAll()

	btnCloseObj, err := b.GetObject("btn_close")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnClose := btnCloseObj.(*gtk.Button)
	btnClose.Connect("clicked", func() { gtk.MainQuit() })

	//окно настроек
	winSettingsObj, err := b.GetObject("window_settings")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	winSettings := winSettingsObj.(*gtk.Window)
	winSettings.Connect("delete-event", func() bool {
		winSettings.Hide()
		return true
	})

	btnCloseSettingsObj, err := b.GetObject("btn_close_settings")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnCloseSettings := btnCloseSettingsObj.(*gtk.Button)
	btnCloseSettings.Connect("clicked", func() { winSettings.Hide() })

	btnSettingsObj, err := b.GetObject("btn_settings")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnSettings := btnSettingsObj.(*gtk.Button)
	btnSettings.Connect("clicked", func() { winSettings.Show() })

	//окно сценариев
	winScriptListObj, err := b.GetObject("window_script_list")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	winScriptList := winScriptListObj.(*gtk.Window)
	winScriptList.Connect("delete-event", func() bool {
		winScriptList.Hide()
		return true
	})

	btnCloseScriptObj, err := b.GetObject("btn_close_script")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnCloseScript := btnCloseScriptObj.(*gtk.Button)
	btnCloseScript.Connect("clicked", func() { winScriptList.Hide() })

	//окно задач
	winTaskListObj, err := b.GetObject("window_task_list")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	winTaskList := winTaskListObj.(*gtk.Window)
	winTaskList.Connect("delete-event", func() bool {
		winTaskList.Hide()
		return true
	})

	btnCloseTaskObj, err := b.GetObject("btn_close_task")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnCloseTask := btnCloseTaskObj.(*gtk.Button)
	btnCloseTask.Connect("clicked", func() { winTaskList.Hide() })

	btnSaveSettingsObj, err := b.GetObject("btn_save_settings")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnSaveSettings := btnSaveSettingsObj.(*gtk.Button)
	btnSaveSettings.Connect("clicked", func() {
		entryNameObj, err := b.GetObject("entry_name")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		entryName := entryNameObj.(*gtk.Entry)
		name, err := entryName.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

		entryAddressObj, err := b.GetObject("entry_address")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

		entryAddress := entryAddressObj.(*gtk.Entry)
		address, err := entryAddress.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

		if setTemp == 0 {
			Conf.Server.Name = name
			Conf.Server.Address = address
		} else if setTemp == 1 {
			Conf.DB.Host = name
			Conf.DB.Port = address
		} else {
			Conf.VMName = name
		}

	})
	//окно резерного копирования БД
	winDataBaseSettingsObj, err := b.GetObject("window_backup_db")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	winDataBaseSettings := winDataBaseSettingsObj.(*gtk.Window)
	winDataBaseSettings.Connect("delete-event", func() bool {
		winDataBaseSettings.Hide()
		return true
	})

	btnDBExitObj, err := b.GetObject("btn_db_exit")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnDBExit := btnDBExitObj.(*gtk.Button)
	btnDBExit.Connect("clicked", func() { winDataBaseSettings.Hide() })

	btnDbOkObj, err := b.GetObject("btn_db_ok")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnDbOk := btnDbOkObj.(*gtk.Button)

	btnForBackupDatabaseObj, err := b.GetObject("btn_backup_db")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnForBackupDatabase := btnForBackupDatabaseObj.(*gtk.Button)
	btnForBackupDatabase.Connect("clicked", func() { winDataBaseSettings.Show() })

	// окно резервного копирования ВМ

	btnBackupVMObj, err := b.GetObject("btn_backup_vm")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnBackupVM := btnBackupVMObj.(*gtk.Button)
	btnBackupVM.Connect("clicked", func() {})
	//окно сжатия и отправки файла
	winCompressFileObj, err := b.GetObject("window_share")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	winCompressFile := winCompressFileObj.(*gtk.Window)
	winCompressFile.Connect("delete-event", func() bool {
		winCompressFile.Hide()
		return true
	})

	entrySharePathNameObj, err := b.GetObject("share_path_name")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	entrySharePathName := entrySharePathNameObj.(*gtk.Entry)

	btnExitShareObj, err := b.GetObject("btn_exit_share")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnExitShare := btnExitShareObj.(*gtk.Button)
	btnExitShare.Connect("clicked", func() { winCompressFile.Hide() })

	//окно настроек
	btnDbOk.Connect("clicked", func() {
		entryDBNameObj, err := b.GetObject("entry_db_name")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		entryDBName := entryDBNameObj.(*gtk.Entry)
		dbName, err := entryDBName.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

		/*entryHostObj, err := b.GetObject("entry_host")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		entryHost := entryHostObj.(*gtk.Entry)

		host, err := entryHost.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

		entryPortObj, err := b.GetObject("entry_port")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		entryPort := entryPortObj.(*gtk.Entry)

		port, err := entryPort.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}*/

		/*entryUsernameObj, err := b.GetObject("entry_username")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		entryUsername := entryUsernameObj.(*gtk.Entry)

		username, err := entryUsername.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

		entryPasswordObj, err := b.GetObject("entry_password")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		entryPassword := entryPasswordObj.(*gtk.Entry)

		password, err := entryPassword.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}*/

		path = copy_db.BackupReindexer(Conf.DB.Host, Conf.DB.Port, dbName)
		winDataBaseSettings.Hide()
		winCompressFile.Show()
		entrySharePathName.SetText(path)
	})

	winBackupVmObj, err := b.GetObject("window_backup_vm")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	winBackupVm := winBackupVmObj.(*gtk.Window)
	winBackupVm.Connect("delete-event", func() bool {
		winBackupVm.Hide()
		return true
	})

	btnBackupVmObj, err := b.GetObject("btn_backup_vm")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnBackupVm := btnBackupVmObj.(*gtk.Button)
	btnBackupVm.Connect("clicked", winBackupVm.Show)

	btnVmDoneObj, err := b.GetObject("btn_vm_ok")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnVmDone := btnVmDoneObj.(*gtk.Button)
	btnVmDone.Connect("clicked", func() {
		entryVmNameObj, err := b.GetObject("entry_vm_name")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		entryVmName := entryVmNameObj.(*gtk.Entry)
		name, err := entryVmName.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		path = copy_db.BackupVM(name)
		winBackupVm.Hide()
		winCompressFile.Show()
		entrySharePathName.SetText(path)
	})

	btnShareCompressedFileObj, err := b.GetObject("btn_share_compressed_file")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	btnShareCompressedFile := btnShareCompressedFileObj.(*gtk.Button)

	btnShareCompressedFile.Connect("clicked", func() {
		sharePathNameObj, err := b.GetObject("share_path_name")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		sharePathName := sharePathNameObj.(*gtk.Entry)

		pathShare, err := sharePathName.GetText()
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

		comboBoxCompressionLevelObj, err := b.GetObject("combo_box_compression_level")
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		comboBoxCompressionLevel := comboBoxCompressionLevelObj.(*gtk.ComboBoxText)
		compLevel := comboBoxCompressionLevel.GetActive()

		/*
			comboBoxProtectionObj, err := b.GetObject("combo_box_protection")
			if err != nil {
				log.Fatal("Ошибка:", err)
			}
			comboBoxProtection := comboBoxProtectionObj.(*gtk.ComboBoxText)
			protectionLevel := comboBoxProtection.GetActive()

			sharePasswordObj, err := b.GetObject("share_password")
			if err != nil {
				log.Fatal("Ошибка:", err)
			}
			sharePassword := sharePasswordObj.(*gtk.Entry)

			password, err := sharePassword.GetText()
			if err != nil {
				log.Fatal("Ошибка:", err)
			}

			numberVirtualMachinesObj, err := b.GetObject("number_virtual_machines")
			if err != nil {
				log.Fatal("Ошибка:", err)
			}
			numberVirtualMachines := numberVirtualMachinesObj.(*gtk.SpinButton)

			numVM, err := numberVirtualMachines.GetText()
			if err != nil {
				log.Fatal("Ошибка:", err)
			}*/

		client_path, client_name := transform.CodeFile(path, pathShare, compLevel-1, "Отсутсвует", "00")
		client.Client(client_path, client_name)
	})

	gtk.Main()
}

func handleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		//test(ctx)
	case "/hello":
		handleHello(ctx)
	default:
		handleNotFound(ctx)
	}
	// Отправляем ответ клиенту
}

// Обработка корневого маршрута "/"
func handleRoot(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Добро пожаловать на главную страницу!")
}

// Обработка маршрута "/hello"
func handleHello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Привет, мир!")
}

// Обработка неизвестных маршрутов
func handleNotFound(ctx *fasthttp.RequestCtx) {
	ctx.Error("Страница не найдена", fasthttp.StatusNotFound)
}

func main() {
	f, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	if err = json.Unmarshal(data, &Conf); err != nil {
		log.Fatal("Ошибка:", err)
	}

	go server.ServerReceiving()
	f.Close()
	gui_pro()

	os.Remove("config/config.json")
	conf, err := os.Create("config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	confJson, err := json.Marshal(Conf)
	if err != nil {
		log.Fatal(err)
	}

	conf.Write(confJson)
	/*file_name := "file.zst"
	file_names := "file.txt"
	// Открываем файл на чтение
	files, err := os.Open(file_names)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer files.Close()
	file, err := os.Create(file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	start := time.Now()
	err = compression.CompressOne(files, file, zstd.EncoderLevelFromZstd(5))
	client.Client(file_name)
	duration := time.Since(start)
	fmt.Println(duration)
	// Nanoseconds как int64
	fmt.Println(duration.Nanoseconds())
	filesf, err := os.Open("test.zst")
	if err != nil {
		panic(err)
	}
	filess, err := os.Create("tests")
	if err != nil {
		panic(err)
	}
	compression.DecompressOne(filesf, filess)
	file.Close()
	files.Close()*/
}
