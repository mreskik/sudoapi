package setup

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// DeviceHostname (2026-10-16, sistem session POS dinamis multi-device) -- hostname
	// komputer POS yang install (hasil gethostname() PHP), dikirim di body GetDataBranch,
	// disimpen ke master_branch_pos_session.device_hostname. Opsional -- endpoint lain yang
	// reuse struct Login ini (GetBranchList, dkk) gak kirim field ini, otomatis kosong.
	DeviceHostname string `json:"device_hostname"`
}

type DataUser struct {
	Id       int    `bun:"id"`
	Username string `bun:"username"`
	Password string `bun:"password"`
}
