package orcl

import (
	"context"
	"strconv"

	// "github.com/alrogozin/lvs_mwr/internal/urq_task"
	"alrogozin/pkg/common"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	// "go.uber.org/zap"

	goora "github.com/sijms/go-ora/v2"
)

var ctx = context.Background()
var dbParams = make(map[string]string, 10)

func Init() {
}

// GetConnection...
func GetConnection(p_db_name string, p_ip_address string, p_port int, p_user string, p_password string) *sqlx.DB {
	dbParams["PORT"] = strconv.Itoa(p_port)
	dbParams["USER"] = p_user
	dbParams["PASSWORD"] = p_password
	dbParams["IP_ADDRESS"] = p_ip_address
	dbParams["DBNAME"] = p_db_name

	dataSourceName := goora.BuildUrl(dbParams["IP_ADDRESS"], p_port, dbParams["DBNAME"], dbParams["USER"], dbParams["PASSWORD"], nil)
	db, err := sqlx.Connect("oracle", dataSourceName)
	if err != nil {
		common.DieOnError("error on connect:", err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		common.DieOnError("error on ping:", err)
		return nil
	}
	// zap.L().Info("Connected to DB " + GetDBName(db))
	return db
}

// GetDBName returnd name of instance from v$database.db_unique_name
func GetDBName(db *sqlx.DB) string {
	var mRetValue string

	sqlString := `Select nvl(db_unique_name, '-') dbname From v$database`
	rows, err := db.Query(sqlString)
	if err != nil {
		common.DieOnError("Can't create query", err)
	}
	for rows.Next() {
		if err = rows.Scan(&mRetValue); err != nil {
			common.DieOnError("On Exec query:", err)
		}
	}

	return mRetValue
}

// PrintVersionInfo...
func PrintVersionInfo(db *sqlx.DB) {
	rows, err := db.Query("SELECT banner FROM v$version")
	if err != nil {
		common.DieOnError("Can't create query", err)
	}
	var banner string
	for rows.Next() {
		if err = rows.Scan(&banner); err != nil {
			common.DieOnError("On Exec query:", err)
		}
		zap.L().Info(banner)
	}
}

/*
// Select id, date_in, date_end, content, tss_name, ustt_name, prs_fio, priority_name, subsys_code, v.cnt_answ
func Get_TL(db *sqlx.DB, marr *[]urq_task.Urq_Task) error {
	var mRet urq_task.Urq_Task
	rows, err := db.Query(common.Param_Oracle.SQL_get_last_tl)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.Scan(&mRet.Id, &mRet.Date_in, &mRet.Date_End, &mRet.Content, &mRet.Tss_name, &mRet.Ustt_name, &mRet.Prs_fio, &mRet.Priority_name, &mRet.Subsys_code, &mRet.Cnt_answers); err != nil {
			return err
		}
		// mRet.Messages = make([]string, 0)
		*marr = append(*marr, mRet)
	}

	return nil
}

// Select cnt_id, date_in, content, num_pp, user_in, type_qr, author_id, author_name, tl_id from urq_vtask_content_mwr v where v.tl_id = :1 order by cnt_id
func Get_qas(db *sqlx.DB, marr *[]urq_task.Urq_qas, p_tl_id int) error {
	var mRet urq_task.Urq_qas
	rows, err := db.Query(common.Param_Oracle.SQL_get_qas, p_tl_id)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.Scan(&mRet.Cnt_id, &mRet.Date_in, &mRet.Content, &mRet.Num_pp, &mRet.User_in, &mRet.Type_qr, &mRet.Author_id, &mRet.Author_name, &mRet.Tl_id); err != nil {
			return err
		}
		*marr = append(*marr, mRet)
	}

	return nil
}
*/
/*
// Select id, chat_abbr, tgn_type, to_char(date_in, 'DD.MM.YY HH24:MI:SS') date_in_dsp, tgn_Grp, decode(get_real_db_name(), 'BILLH', 'Prod', 'SDAILY', 'Daily', get_real_db_name()) db_name from tgn_hdr where date_sent is null order by id
func FillTgMessage(db *sqlx.DB, marr *[]common.TGMessage_Type) error {
	var mRet common.TGMessage_Type
	rows, err := db.Query(common.Param_Oracle.SQL4Process)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.Scan(&mRet.Id, &mRet.Chat_Abbr, &mRet.Tgn_Type, &mRet.Date_in, &mRet.Tgn_Grp, &mRet.Db_Name); err != nil {
			return err
		}
		mRet.Messages = make([]string, 0)
		*marr = append(*marr, mRet)
	}

	return nil
}

// Select msg_text from tgn_data where hdr_id = :1 order by num_order
func FillMessages(db *sqlx.DB, hdr_id int, msgs *[]string) error {
	var mRet string
	rows, err := db.Query(common.Param_Oracle.SQL_get_message, hdr_id)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.Scan(&mRet); err != nil {
			return err
		}
		*msgs = append(*msgs, mRet)
	}

	return nil
}

// Update tgn_hdr set date_sent = sysdate where id = :1
func Set_Done(db *sqlx.DB, hdr_id int) error {
	tx, _ := db.Begin()
	_, err := tx.Exec(common.Param_Oracle.SQL_set_done, hdr_id)
	if err != nil {
		return err
	}
	_ = tx.Commit()
	return nil
}
*/

func SetAppInfo(db *sqlx.DB, p_message string) error {
	_, err := db.Exec(common.Param_Oracle.SQL_set_appinfo, p_message)
	if err != nil {
		zap.L().Sugar().Errorf("%v %s", err, "Set_AppInfo "+p_message)
		return err
	}
	return nil
}

/*
// "SQL_get_unprocessed_records" : "Select id, sbs_code, gpzu_id, rd_id, ops_id_inp, ldr_id_doc, ldr_id_sign, status_his_id from gpze_data where gpze_state = 'INIT' order by id"
func SQL_get_unprocessed_records(db *sqlx.DB) ([]gpze.Gpze, error) {
	rows, err := db.Query(common.Param_Oracle.SQL4Process)
	common.DieOnError("Запрос gpze_data ", err)
	var gpzeArr []gpze.Gpze = make([]gpze.Gpze, 5)
	defer func() {
		err = rows.Close()
		common.DieOnError("Cant close dataset ", err)
	}()
	gpze := gpze.Gpze{}
	for rows.Next() {
		err = rows.Scan(&gpze.Id, &gpze.Sbs_Code, &gpze.Gpzu_Id, &gpze.Rd_Id, &gpze.Ops_Id_Inp, &gpze.Ldr_Id_Doc, &gpze.Ldr_Id_Sign, &gpze.Status_His_Id)
		// common.DieOnError("Scan ", err)
		if err != nil {
			return nil, err
		}
		gpzeArr = append(gpzeArr, gpze)
	}

	return gpzeArr, nil
}

func SQL_get_file(db *sqlx.DB, p_gpze *gpze.Gpze, p_type string) error {
	var p_id int
	if p_type == "DOC" {
		p_id = p_gpze.Ldr_Id_Doc
	} else {
		p_id = p_gpze.Ldr_Id_Sign
	}

	rows, err := db.Query(common.Param_Oracle.SQL_GetFile, p_id)
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		if p_type == "DOC" {
			err = rows.Scan(&p_gpze.Doc_File_name, &p_gpze.Doc_blob)
		} else {
			err = rows.Scan(&p_gpze.Sign_File_name, &p_gpze.Sign_blob)
		}
	}

	return nil
}

// Begin gpze_mng_pkg.Load_File2Ldr(p_gpze_id => :1, p_file_name => :2, p_type => :3, p_data => :4); End;
func SaveUnzipData(db *sqlx.DB, p_gpze_id int, p_file_name string, p_blob_data []byte, p_type string) error {
	var gb go_ora.Blob
	gb.Scan(p_blob_data)
	_, err := db.Exec(common.Param_Oracle.SQL_LoadUnzipFiles, p_gpze_id, p_file_name, p_type, gb)
	if err != nil {
		zap.L().Sugar().Errorf("%v %s", err, p_type)
		return err
	}

	return nil
}

// "SQL_save_verdict": "Begin gpze_mng_pkg.save_verdict(p_gpze_id => :1, p_verdict => :2, p_message => :3); End;"
func Save_Verdict(db *sqlx.DB, p_gpze_id int, p_verdict string, p_message string) error {
	_, err := db.Exec(common.Param_Oracle.SQL_save_verdict, p_gpze_id, p_verdict, p_message)
	if err != nil {
		zap.L().Sugar().Errorf("%v %s", err, "SAVE_VERDICT "+p_verdict)
		return err
	}

	return nil
}

// "SQL_save_cert_data": "Begin GPZE_MNG_PKG.Save_Cert_Data(p_gpze_id => :1, p_signing_date => :2, p_serial_number => :3, p_thumbprint => :4, p_snils => :5, p_ogrn => :6, p_inn => :7, p_fio => :8, p_location => :9); End;"
func Save_Cert_Data(db *sqlx.DB, p_gpze *gpze.Gpze) error {
	_, err := db.Exec(common.Param_Oracle.SQL_save_cert_data, p_gpze.Id, p_gpze.SigningDateString, p_gpze.SerialNumber, p_gpze.Thumbprint, p_gpze.Snils, p_gpze.Ogrn, p_gpze.Inn, p_gpze.Fio, p_gpze.Location)
	if err != nil {
		zap.L().Sugar().Errorf("%v %s", err, "SAVE_CERT")
		return err
	}

	return nil
}

// "SQL_save_error": "Begin SaveError(p_gpze_id => :1, p_type => :2, p_message => :3); End;"
func Save_Error(db *sqlx.DB, p_gpze_id int, p_type string, p_message string) error {
	_, err := db.Exec(common.Param_Oracle.SQL_save_error, p_gpze_id, p_message, p_message)
	if err != nil {
		panic(err)
	}

	return nil
}

// "SQL_write_tgn":  "Declare mBuff number; Begin mBuff := tgn_pkg.reg_simple_tgn(p_abbr => 'LeavesAlert', p_type => 'ERROR', p_grp => 'GPZE.Error_Inner', p_txt => :1, p_ref_id => :2); End;"
func Write_tgn(db *sqlx.DB, p_gpze_id int, p_message string, p_type string) error {
	_, err := db.Exec(common.Param_Oracle.SQL_write_tgn, p_type, p_message, p_gpze_id)
	if err != nil {
		panic(err)
	}

	return nil
}

func Reset_Gpze(db *sqlx.DB, p_gpze_id int) error {
	_, err := db.Exec(common.Param_Oracle.SQL_reset_gpze, p_gpze_id)
	if err != nil {
		panic(err)
	}

	return nil
}
*/
/*
func InitLoad(db *sqlx.DB, menu_name string) (int, error) {
	const pl_block_init = "Begin :1 := omenu_src_pkg.Init_Load(:2); end;"
	var hdr_id int
	_, err := db.Exec(pl_block_init, &hdr_id, menu_name)
	if err != nil {
		panic(err)
	}

	return hdr_id, nil
}
*/
/*
func InsMenuItem(db *sqlx.DB, unit manageMMT.MenuItem, hdr_id int) error {

	const pl_block = "Begin omenu_src_pkg.Load_Unit(:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13); end;"
	_, err := db.Exec(pl_block, &unit.Name, &unit.Label, &unit.MenuType, &unit.CommandType, &unit.CommandType, &unit.Visible, &unit.Dsp_wo_privs, &unit.MenuItemCode, &unit.Id, &unit.Prnt_id, &unit.StartPos, &unit.EndPos, &hdr_id)
	if err != nil {
		panic(err)
	}

	const pl_block_r = "Begin omenu_src_pkg.Load_Roles(:1, :2, :3); end;"
	for _, r := range unit.MenuRoles {
		_, err = db.Exec(pl_block_r, &unit.Id, r, &hdr_id)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
*/
/*
// DoDDL ...
func DoDDL(db *sqlx.DB, ddlString string) {
	_, err := db.Exec(ddlString)
	if err != nil {
		common.DieOnError("error:", err)
	}

}
*/
