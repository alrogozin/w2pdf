package common

// Select id, chat_abbr, tgn_type, to_char(date_in, 'DD.MM.YY HH24:MI:SS') date_in_dsp, tgn_Grp, decode(get_real_db_name(), 'BILLH', 'Prod', 'SDAILY', 'Daily', get_real_db_name()) db_name from tgn_hdr where date_sent is null order by id
type TGMessage_Type struct {
	Id        int
	Chat_Abbr string
	Tgn_Type  string
	Date_in   string
	Tgn_Grp   string
	Db_Name   string
	Messages  []string
}
