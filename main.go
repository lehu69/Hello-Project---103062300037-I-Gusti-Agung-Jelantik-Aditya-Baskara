package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type Barang struct {
	ID       int
	Nama     string
	Jumlah   int
	Kategori string
	Waktu    time.Time
}

type BarangDihapus struct {
	ID         int
	Nama       string
	Jumlah     int
	Kategori   string
	Waktu      time.Time
	HapusWaktu time.Time
}

// Save data to a file
func saveDataToFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

// Format data to text
func formatDataToText(dataBarang []Barang) string {
	formattedData := "Daftar Barang:\n"
	for _, barang := range dataBarang {
		formattedData += fmt.Sprintf("ID: %d, Nama: %s, Jumlah: %d, Kategori: %s, Waktu: %s\n",
			barang.ID, barang.Nama, barang.Jumlah, barang.Kategori, barang.Waktu.Format("2006-01-02 15:04:05"))
	}
	return formattedData
}

// menyimpan data barang yang di hapus
func formatDeletedDataToText(dataBarangDihapus []BarangDihapus) string {
	formattedData := "Daftar Barang Dihapus:\n"
	for _, barang := range dataBarangDihapus {
		formattedData += fmt.Sprintf("ID: %d, Nama: %s, Jumlah: %d, Kategori: %s, Waktu Dihapus: %s\n",
			barang.ID, barang.Nama, barang.Jumlah, barang.Kategori, barang.HapusWaktu.Format("2006-01-02 15:04:05"))
	}
	return formattedData
}

// menyimpan data
func saveData(dataBarang []Barang, dataBarangDihapus []BarangDihapus) {
	fileGob, err := os.Create("data.gob")
	if err != nil {
		fmt.Println("Error saving data:", err)
		return
	}
	defer fileGob.Close()

	encoder := gob.NewEncoder(fileGob)
	err = encoder.Encode(dataBarang)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

	fmt.Println("Data berhasil disimpan dalam format GOB.")

	// Save data to text file
	dataText := formatDataToText(dataBarang)
	err = saveDataToFile("data.txt", dataText)
	if err != nil {
		fmt.Println("Error saving data to text file:", err)
		return
	}
	fmt.Println("Data berhasil disimpan dalam format teks.")

	// Save deleted data to text file
	deletedDataText := formatDeletedDataToText(dataBarangDihapus)
	err = saveDataToFile("deleted_data.txt", deletedDataText)
	if err != nil {
		fmt.Println("Error saving deleted data to text file:", err)
		return
	}
	fmt.Println("Data barang dihapus berhasil disimpan dalam format teks.")
}

// load data yang sudah di masukan sebelumnya
func loadData() ([]Barang, []BarangDihapus, error) {
	var dataBarang []Barang
	var dataBarangDihapus []BarangDihapus

	file, err := os.Open("data.gob")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File data tidak ditemukan, memulai dengan data kosong.")
			return dataBarang, dataBarangDihapus, nil
		}
		return nil, nil, fmt.Errorf("Error loading data: %v", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&dataBarang)
	if err != nil {
		return nil, nil, fmt.Errorf("Error decoding dataBarang: %v", err)
	}

	fmt.Println("Data berhasil dimuat.")
	return dataBarang, dataBarangDihapus, nil
}

// membuat tabel agar data terlihat rapi
func printTable(dataBarang []Barang) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No", "Nama Barang", "ID", "Jumlah", "Kategori", "Waktu"})

	for i, barang := range dataBarang {
		table.Append([]string{
			fmt.Sprintf("%d", i+1),
			barang.Nama,
			fmt.Sprintf("%d", barang.ID),
			fmt.Sprintf("%d", barang.Jumlah),
			barang.Kategori,
			barang.Waktu.Format("2006-01-02 15:04:05"),
		})
	}
	table.Render()
}

// tabel untuk data yang sudah terhapus
func printDeletedTable(dataBarangDihapus []BarangDihapus) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No", "Nama Barang", "ID", "Jumlah", "Kategori", "Waktu Dihapus"})

	for i, barang := range dataBarangDihapus {
		table.Append([]string{
			fmt.Sprintf("%d", i+1),
			barang.Nama,
			fmt.Sprintf("%d", barang.ID),
			fmt.Sprintf("%d", barang.Jumlah),
			barang.Kategori,
			barang.HapusWaktu.Format("2006-01-02 15:04:05"),
		})
	}
	table.Render()
}

func cariIndeksBarang(dataBarang []Barang, namaBarang string) int {
	for i, barang := range dataBarang {
		if barang.Nama == namaBarang {
			return i
		}
	}
	return -1
}

// mengubah nama barang
func renameBarang(dataBarang []Barang, namaLama string, namaBaru string) {
	indeks := cariIndeksBarang(dataBarang, namaLama)
	if indeks != -1 {
		dataBarang[indeks].Nama = namaBaru
		color.Green("Nama barang berhasil diganti.")
	} else {
		color.Red("Barang tidak ditemukan")
	}
}

// menghapus data barang
func deleteBarang(dataBarang []Barang, dataBarangDihapus []BarangDihapus, namaBarang string, count int) ([]Barang, []BarangDihapus, int) {
	indeks := cariIndeksBarang(dataBarang, namaBarang)
	if indeks != -1 {
		barangDihapus := BarangDihapus{
			ID:         dataBarang[indeks].ID,
			Nama:       dataBarang[indeks].Nama,
			Jumlah:     dataBarang[indeks].Jumlah,
			Kategori:   dataBarang[indeks].Kategori,
			Waktu:      dataBarang[indeks].Waktu,
			HapusWaktu: time.Now(),
		}
		dataBarangDihapus = append(dataBarangDihapus, barangDihapus)
		dataBarang = append(dataBarang[:indeks], dataBarang[indeks+1:]...)
		count--
		color.Green("Barang berhasil dihapus pada waktu %s.", time.Now().Format("2006-01-02 15:04:05"))
	} else {
		color.Red("Barang tidak ditemukan")
	}
	return dataBarang, dataBarangDihapus, count
}

// mencari data barang berdasarkan kategori
func searchBarangByKategori(dataBarang []Barang, kategoriBarang string) {
	color.Cyan("Hasil pencarian:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No", "Nama Barang", "ID", "Jumlah", "Kategori", "Waktu"})
	found := false
	for i, barang := range dataBarang {
		if barang.Kategori == kategoriBarang {
			table.Append([]string{
				fmt.Sprintf("%d", i+1),
				barang.Nama,
				fmt.Sprintf("%d", barang.ID),
				fmt.Sprintf("%d", barang.Jumlah),
				barang.Kategori,
				barang.Waktu.Format("2006-01-02 15:04:05"),
			})
			found = true
		}
	}
	if !found {
		color.Red("Barang tidak ditemukan")
	} else {
		table.Render()
	}
}

func updateJumlahBarang(dataBarang []Barang, namaBarang string, jumlahTambah int) {
	indeks := cariIndeksBarang(dataBarang, namaBarang)
	if indeks != -1 {
		dataBarang[indeks].Jumlah += jumlahTambah
		dataBarang[indeks].Waktu = time.Now() // Update the time to the current time
		color.Green(fmt.Sprintf("Jumlah %s berhasil ditambahkan. Jumlah sekarang: %d. Waktu diperbarui ke %s.", namaBarang, dataBarang[indeks].Jumlah, dataBarang[indeks].Waktu.Format("2006-01-02 15:04:05")))
	} else {
		color.Red("Barang tidak ditemukan")
	}
}

func kurangiJumlahBarang(dataBarang []Barang, dataBarangDihapus []BarangDihapus, namaBarang string, jumlahKurang int, count int) ([]Barang, []BarangDihapus, int) {
	indeks := cariIndeksBarang(dataBarang, namaBarang)
	if indeks != -1 {
		if dataBarang[indeks].Jumlah > jumlahKurang {
			dataBarang[indeks].Jumlah -= jumlahKurang
			dataBarang[indeks].Waktu = time.Now() // Update the time to the current time
			color.Green(fmt.Sprintf("Jumlah %s berhasil dikurangi. Jumlah sekarang: %d. Waktu diperbarui ke %s.", namaBarang, dataBarang[indeks].Jumlah, dataBarang[indeks].Waktu.Format("2006-01-02 15:04:05")))
		} else {
			barangDihapus := BarangDihapus{
				ID:         dataBarang[indeks].ID,
				Nama:       dataBarang[indeks].Nama,
				Jumlah:     dataBarang[indeks].Jumlah,
				Kategori:   dataBarang[indeks].Kategori,
				Waktu:      dataBarang[indeks].Waktu,
				HapusWaktu: time.Now(),
			}
			dataBarangDihapus = append(dataBarangDihapus, barangDihapus)
			dataBarang = append(dataBarang[:indeks], dataBarang[indeks+1:]...)
			count--
			color.Green("Barang %s berhasil dihapus karena jumlah barang menjadi nol atau kurang. Dihapus pada waktu %s.", namaBarang, time.Now().Format("2006-01-02 15:04:05"))
		}
	} else {
		color.Red("Barang tidak ditemukan")
	}
	return dataBarang, dataBarangDihapus, count
}

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Platform tidak didukung.")
	}
}

func sortBynama(dataBarang []Barang) {
	n := len(dataBarang)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if dataBarang[j].Nama < dataBarang[minIdx].Nama {
				minIdx = j
			}
		}
		dataBarang[i], dataBarang[minIdx] = dataBarang[minIdx], dataBarang[i]
	}
	color.Green("Data berhasil diurutkan berdasarkan nama.")
	printTable(dataBarang)
}

func sortByjumlah(dataBarang []Barang) {
	n := len(dataBarang)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if dataBarang[j].Jumlah < dataBarang[minIdx].Jumlah {
				minIdx = j
			}
		}
		dataBarang[i], dataBarang[minIdx] = dataBarang[minIdx], dataBarang[i]
	}
	color.Green("Data berhasil diurutkan berdasarkan jumlah.")
	printTable(dataBarang)
}

func main() {
	var namaBarang string

	dataBarang, dataBarangDihapus, err := loadData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	count := len(dataBarang)

	for {
		color.Cyan("\nPilih operasi yang ingin Anda lakukan:")
		fmt.Println("+----+-----------------------+")
		fmt.Println("| 1. | Tampilkan data barang |")
		fmt.Println("+----+-----------------------+")
		fmt.Println("| 2. | Cari Barang (Kategori)|")
		fmt.Println("+----+-----------------------+")
		fmt.Println("| 3. | Urutkan Barang        |")
		fmt.Println("+----+-----------------------+")
		fmt.Println("| 4. | Update barang         |")
		fmt.Println("+----+-----------------------+")
		fmt.Println("| 5. | Data barang Dihapus   |")
		fmt.Println("+----+-----------------------+")
		fmt.Println("| 6. | Simpan ke file        |")
		fmt.Println("+----+-----------------------+")
		fmt.Println("| 0. | Keluar                |")
		fmt.Println("+----+-----------------------+")

		var pilihan int
		fmt.Print("\nMasukkan pilihan: ")
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			printTable(dataBarang)
		case 2:
			fmt.Print("Masukkan kategori barang yang ingin Anda cari: ")
			var kategoriBarang string
			fmt.Scan(&kategoriBarang)
			searchBarangByKategori(dataBarang, kategoriBarang)
		case 3:
			color.Cyan("\nPilih operasi yang ingin Anda lakukan:")
			fmt.Println("+----+------------------------------------------+")
			fmt.Println("| 1. | Urutkan barang berdasarkan nama          |")
			fmt.Println("+----+------------------------------------------+")
			fmt.Println("| 2. | urutkan barang berdasarkan jumlah barang |")
			fmt.Println("+----+------------------------------------------+")

			var urutBarang int
			fmt.Print("\nMasukkan pilihan: ")
			fmt.Scan(&urutBarang)
			switch urutBarang {
			case 1:
				sortBynama(dataBarang)
			case 2:
				sortByjumlah(dataBarang)
			}
		case 4:
			color.Cyan("\nPilih update yang ingin Anda lakukan:")
			fmt.Println("+----+----------------------+")
			fmt.Println("| 1. | Tambah barang baru   |")
			fmt.Println("+----+----------------------+")
			fmt.Println("| 2. |Tambah jumlah barang  |")
			fmt.Println("+----+----------------------+")
			fmt.Println("| 3. | Kurangi jumlah barang|")
			fmt.Println("+----+----------------------+")
			fmt.Println("| 4. | Hapus data barang    |")
			fmt.Println("+----+----------------------+")
			fmt.Println("| 5. | Ganti nama barang    |")
			fmt.Println("+----+----------------------+")
			fmt.Print("\nMasukkan pilihan: ")
			var updatePilihan int
			fmt.Scan(&updatePilihan)
			switch updatePilihan {
			case 1:
				if count < 100 {
					var namaBaru, kategoriBaru string
					var jumlahBaru int
					fmt.Print("Masukkan nama barang baru: ")
					fmt.Scan(&namaBaru)
					fmt.Print("Masukkan jumlah barang baru: ")
					fmt.Scan(&jumlahBaru)
					fmt.Print("Masukkan kategori barang baru: ")
					fmt.Scan(&kategoriBaru)
					dataBarang = append(dataBarang, Barang{
						ID:       count + 1,
						Nama:     namaBaru,
						Jumlah:   jumlahBaru,
						Kategori: kategoriBaru,
						Waktu:    time.Now(),
					})
					count++
					color.Green("Barang berhasil ditambahkan.")
				} else {
					color.Red("Data barang penuh, tidak bisa menambah barang lagi.")
				}
			case 2:
				fmt.Print("Masukkan nama barang yang ingin Anda tambahkan jumlahnya: ")
				fmt.Scan(&namaBarang)
				fmt.Print("Masukkan jumlah yang ingin ditambahkan: ")
				var jumlahTambah int
				fmt.Scan(&jumlahTambah)
				updateJumlahBarang(dataBarang, namaBarang, jumlahTambah)
			case 3:
				fmt.Print("Masukkan nama barang yang ingin Anda kurangi jumlahnya: ")
				fmt.Scan(&namaBarang)
				fmt.Print("Masukkan jumlah yang ingin dikurangi: ")
				var jumlahKurang int
				fmt.Scan(&jumlahKurang)
				dataBarang, dataBarangDihapus, count = kurangiJumlahBarang(dataBarang, dataBarangDihapus, namaBarang, jumlahKurang, count)
			case 4:
				fmt.Print("Masukkan nama barang yang ingin Anda hapus: ")
				fmt.Scan(&namaBarang)
				dataBarang, dataBarangDihapus, count = deleteBarang(dataBarang, dataBarangDihapus, namaBarang, count)
			case 5:
				fmt.Print("Masukkan nama barang yang ingin Anda ganti: ")
				var namaLama, namaBaru string
				fmt.Scan(&namaLama)
				fmt.Print("Masukkan nama baru untuk barang ", namaLama, ": ")
				fmt.Scan(&namaBaru)
				renameBarang(dataBarang, namaLama, namaBaru)
			default:
				color.Red("Pilihan tidak valid")
			}
		case 5:
			printDeletedTable(dataBarangDihapus)
		case 6:
			clearScreen()
			saveData(dataBarang, dataBarangDihapus)
		case 0:
			return
		default:
			color.Red("Pilihan tidak valid")
		}
	}
}
