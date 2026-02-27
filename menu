package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	fmt.Println("----- MENU PRINCIPAL -----")
	fmt.Println("A : Analyse sur fichier courant")
	fmt.Println("B : Analyse multi-fichiers")
	fmt.Println("C : Analyser une page Wikipédia")
	fmt.Print("Votre choix : ")

	var choixPrincipal string
	fmt.Scanln(&choixPrincipal)

	if choixPrincipal == "A" {
		menuFichierSimple()
	} else if choixPrincipal == "B" {
		menuMultiFichiers()
	} else {
		fmt.Println("Ecris bien")
	}
}

func menuFichierSimple() {

	fichierConfig, _ := os.Open("config.txt")
	defer fichierConfig.Close()

	var fichierParDefaut string
	var dossierSortie string

	scannerConfig := bufio.NewScanner(fichierConfig)

	for scannerConfig.Scan() {
		ligne := strings.TrimSpace(scannerConfig.Text())
		parties := strings.Split(ligne, "=")

		if len(parties) == 2 {
			if parties[0] == "default_file" {
				fichierParDefaut = parties[1]
			}
			if parties[0] == "out_dir" {
				dossierSortie = parties[1]
			}
		}
	}

	for {

		fmt.Println("------------- MENU -------------")
		fmt.Println("1) Infos fichier")
		fmt.Println("2) Statistiques mots")
		fmt.Println("3) Compter lignes avec mot clé")
		fmt.Println("4) Filtrer lignes avec mot clé")
		fmt.Println("5) Filtrer lignes sans mot clé")
		fmt.Println("6) Head")
		fmt.Println("7) Tail")
		fmt.Println("0) Quitter")

		fmt.Print("Choix : ")

		var choix string
		fmt.Scanln(&choix)

		if choix == "0" {
			break
		}

		fmt.Print("Chemin fichier : ")
		var chemin string
		fmt.Scanln(&chemin)

		if chemin == "" {
			chemin = fichierParDefaut
		}

		lignes := lireFichier(chemin)

		switch choix {

		case "1":
			info, _ := os.Stat(chemin)
			fmt.Println("Taille :", info.Size())
			fmt.Println("Modifié :", info.ModTime().Format(time.RFC822))
			fmt.Println("Nb lignes :", len(lignes))

		case "2":
			nombreMots := 0
			for i := 0; i < len(lignes); i++ {
				mots := strings.Fields(lignes[i])
				nombreMots = nombreMots + len(mots)
			}
			fmt.Println("Nb de mots :", nombreMots)

		case "3":
			fmt.Print("Mot clé : ")
			var mot string
			fmt.Scanln(&mot)
			compteur := 0
			for i := 0; i < len(lignes); i++ {
				if strings.Contains(lignes[i], mot) {
					compteur++
				}
			}
			fmt.Println("Nb lignes :", compteur)

		case "4":
			fmt.Print("Mot clé : ")
			var mot string
			fmt.Scanln(&mot)
			var resultat []string
			for i := 0; i < len(lignes); i++ {
				if strings.Contains(lignes[i], mot) {
					resultat = append(resultat, lignes[i])
				}
			}
			os.MkdirAll(dossierSortie, os.ModePerm)
			ecrireFichier(dossierSortie+"/filtered.txt", resultat)

		case "5":
			fmt.Print("Mot clé : ")
			var mot string
			fmt.Scanln(&mot)
			var resultat []string
			for i := 0; i < len(lignes); i++ {
				if !strings.Contains(lignes[i], mot) {
					resultat = append(resultat, lignes[i])
				}
			}
			os.MkdirAll(dossierSortie, os.ModePerm)
			ecrireFichier(dossierSortie+"/filtered_not.txt", resultat)

		case "6":
			fmt.Print("Nombre : ")
			var nombre int
			fmt.Scanln(&nombre)
			if nombre > len(lignes) {
				nombre = len(lignes)
			}
			os.MkdirAll(dossierSortie, os.ModePerm)
			ecrireFichier(dossierSortie+"/head.txt", lignes[:nombre])

		case "7":
			fmt.Print("Nombre : ")
			var nombre int
			fmt.Scanln(&nombre)
			if nombre > len(lignes) {
				nombre = len(lignes)
			}
			debut := len(lignes) - nombre
			os.MkdirAll(dossierSortie, os.ModePerm)
			ecrireFichier(dossierSortie+"/tail.txt", lignes[debut:])
		}
	}
}

func menuMultiFichiers() {

	for {

		fmt.Println("------ MENU ------")
		fmt.Println("1) Batch analyser tous les .txt")
		fmt.Println("2) Rapport global")
		fmt.Println("3) Indexation")
		fmt.Println("4) Fusion")
		fmt.Println("0) Retour")

		fmt.Print("Choix : ")

		var choix string
		fmt.Scanln(&choix)

		if choix == "0" {
			break
		}

		fmt.Print("Nom du répertoire : ")
		var dossier string
		fmt.Scanln(&dossier)

		fichiers, _ := os.ReadDir(dossier)

		var listeTxt []string

		for i := 0; i < len(fichiers); i++ {
			nom := fichiers[i].Name()
			if strings.HasSuffix(nom, ".txt") {
				listeTxt = append(listeTxt, dossier+"/"+nom)
			}
		}

		os.MkdirAll("out", os.ModePerm)

		if choix == "1" {
			fmt.Println("Batch terminé :", len(listeTxt))
		}

		if choix == "2" {

			var rapport []string
			totalLignes := 0

			for i := 0; i < len(listeTxt); i++ {
				lignes := lireFichier(listeTxt[i])
				totalLignes = totalLignes + len(lignes)
			}

			rapport = append(rapport, "Nb fichiers : "+strconv.Itoa(len(listeTxt)))
			rapport = append(rapport, "Total lignes : "+strconv.Itoa(totalLignes))

			ecrireFichier("out/report.txt", rapport)
		}

		if choix == "3" {

			var index []string

			for i := 0; i < len(listeTxt); i++ {
				info, _ := os.Stat(listeTxt[i])
				ligne := listeTxt[i] + " | " +
					strconv.FormatInt(info.Size(), 10) + " octets | " +
					info.ModTime().Format(time.RFC822)
				index = append(index, ligne)
			}

			ecrireFichier("out/index.txt", index)
		}

		if choix == "4" {

			var fusion []string

			for i := 0; i < len(listeTxt); i++ {
				lignes := lireFichier(listeTxt[i])
				for j := 0; j < len(lignes); j++ {
					fusion = append(fusion, lignes[j])
				}
			}

			ecrireFichier("out/merged.txt", fusion)
		}
	}
}

func lireFichier(chemin string) []string {

	fichier, erreur := os.Open(chemin)
	if erreur != nil {
		return []string{}
	}
	defer fichier.Close()

	var lignes []string

	scanner := bufio.NewScanner(fichier)
	for scanner.Scan() {
		lignes = append(lignes, scanner.Text())
	}

	return lignes
}

func ecrireFichier(chemin string, lignes []string) {

	fichier, _ := os.Create(chemin)
	defer fichier.Close()

	for i := 0; i < len(lignes); i++ {
		fichier.WriteString(lignes[i] + "\n")
	}
}

func menuWikipedia() {

	fmt.Print("Nom article : ")
	var article string
	fmt.Scanln(&article)

	url := "https://fr.wikipedia.org/wiki/" + article

	reponse, erreur := http.Get(url)
	if erreur != nil {
		fmt.Println("Le dl ne fonctionne pas")
		return
	}
	defer reponse.Body.Close()

	document, erreur := goquery.NewDocumentFromReader(reponse.Body)
	if erreur != nil {
		fmt.Println("Y'a pas")
		return
	}

	var paragraphes []string

	document.Find("p").Each(func(index int, element *goquery.Selection) {
		texte := element.Text()
		if strings.TrimSpace(texte) != "" {
			paragraphes = append(paragraphes, texte)
		}
	})

	totalMots := 0
	for i := 0; i < len(paragraphes); i++ {
		mots := strings.Fields(paragraphes[i])
		totalMots = totalMots + len(mots)
	}

	nombreLignes := len(paragraphes)

	fmt.Println("Nb de paragraphes :", nombreLignes)
	fmt.Println("Nb total de mots :", totalMots)

	os.MkdirAll("out", os.ModePerm)

	nomFichier := "out/wiki_" + article + ".txt"

	var contenu []string
	contenu = append(contenu, "Article : "+article)
	contenu = append(contenu, "Nb paragraphes : "+strconv.Itoa(nombreLignes))
	contenu = append(contenu, "Nb mots : "+strconv.Itoa(totalMots))
	contenu = append(contenu, "Contenu : ")

	for i := 0; i < len(paragraphes); i++ {
		contenu = append(contenu, paragraphes[i])
	}

	ecrireFichier(nomFichier, contenu)

	fmt.Println("Fichier généré :", nomFichier)
}
