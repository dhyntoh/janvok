# ZIVPN Installer

CLI untuk mengelola aktivasi dan menjalankan installer ZIVPN + Xray (VMess/VLESS/Trojan) + Hysteria2, termasuk admin bot Telegram untuk manajemen akun.

## Build

```bash
go build -o zivpn-installer .
```

## Cara Pakai (Ringkas)

```bash
./zivpn-installer mint --months 1
export ACTIVATION_TOKEN="token_dari_mint"
./zivpn-installer install
```

## Perintah CLI

```bash
zivpn-installer <command> [options]

Commands:
  mint     Create activation token (1 month or 1 year)
  bot      Run Telegram admin bot for account management
  install  Run installer with activation token
  update   Pull latest repository changes and rebuild
```

## Aktivasi Token

Membuat token aktivasi untuk masa berlaku tertentu.

```bash
./zivpn-installer mint --months 1
./zivpn-installer mint --years 1
```

Opsional menentukan label plan:

```bash
./zivpn-installer mint --months 6 --plan "6-months"
```

Token akan disimpan di `INSTALLER_HOME` (default: `/etc/zivpn-installer` untuk root, atau `~/.zivpn-installer` untuk user non-root) pada file `tokens.json`.

## Instalasi

Jalankan installer dengan token aktivasi:

```bash
export ACTIVATION_TOKEN="token_dari_mint"
./zivpn-installer install
```

Atau masukkan token saat diminta. Setelah token valid, installer akan meminta input interaktif untuk:

- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_ADMIN_IDS`
- domain Xray (`XRAY_DOMAIN`)

Konfigurasi akan disimpan di `INSTALLER_HOME/config.json` dan digunakan kembali.

### Port Default

- Xray TCP: **443**
- Hysteria2 UDP: **443**
- ZIVPN UDP: **5667**

### Dry Run

Untuk melihat langkah instalasi tanpa mengeksekusi:

```bash
./zivpn-installer install --dry-run
```

## Update Otomatis Repo

Jalankan update otomatis dari repo Git dan rebuild binary:

```bash
./zivpn-installer update
```

Jika repo tidak berada di working directory, set lokasi repo:

```bash
export INSTALLER_REPO="/path/to/repo"
./zivpn-installer update
```

## Telegram Admin Bot

Set environment berikut sebelum menjalankan bot:

```bash
export TELEGRAM_BOT_TOKEN="bot_token"
export TELEGRAM_ADMIN_IDS="123456789,987654321"
```

Atau biarkan kosong jika sudah tersimpan di `INSTALLER_HOME/config.json`.

Jalankan bot:

```bash
./zivpn-installer bot
```

### Systemd Bot Service

Saat instalasi, bot otomatis dibuat sebagai service systemd dan dijalankan:

```
systemctl status zivpn-installer-bot.service
```

### Perintah Admin

Gunakan perintah atau menu:

```
/create
/delete <username>
/setexp <username> <YYYY-MM-DD>
/setlimit <username> <limit>
/list
/status
/restart <zivpn|xray|hysteria2|bot>
/backup
/restore [file]
/cancel
```

Flow `/create` bersifat interaktif: bot akan menanyakan username, masa aktif (hari), lalu pilihan protocol dan mengirim link akun seperti `zi://`, `vless://`, dll.

Data akun akan disimpan di `accounts.json` pada `INSTALLER_HOME`.

## Variabel Lingkungan

| Variable | Keterangan |
| --- | --- |
| `ACTIVATION_TOKEN` | Token aktivasi untuk instalasi. |
| `INSTALLER_HOME` | Lokasi penyimpanan token & akun. |
| `INSTALLER_REPO` | Lokasi repo Git untuk perintah update. |
| `XRAY_DOMAIN` | Domain untuk Xray (disimpan ke config saat instalasi). |
| `TELEGRAM_BOT_TOKEN` | Token bot Telegram. |
| `TELEGRAM_ADMIN_IDS` | Daftar ID admin (pisahkan dengan koma). |
| `TELEGRAM_CHAT_ID` | Chat ID untuk notifikasi token/instalasi. |
