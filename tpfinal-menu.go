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
		fmt.Println("En majuscule :)")
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

	var articles []string
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Nom article : ")
		article, _ := reader.ReadString('\n')
		article = strings.TrimSpace(article)
		article = strings.ReplaceAll(article, " ", "_")

		if article == "" {
			fmt.Println("Nom invalide")
			continue
		}

		articles = append(articles, article)

		fmt.Print("Analyser un autre article ? (oui/non) : ")
		choix, _ := reader.ReadString('\n')
		choix = strings.TrimSpace(choix)

		if choix != "oui" {
			break
		}
	}

	for i := 0; i < len(articles); i++ {

		article := articles[i]
		url := "https://fr.wikipedia.org/wiki/" + article
		fmt.Println("Téléchargement :", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Erreur requête :", err)
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")

		client := &http.Client{}
		reponse, err := client.Do(req)
		if err != nil {
			fmt.Println("Erreur HTTP :", err)
			continue
		}

		if reponse.StatusCode != 200 {
			fmt.Println("Article non trouvé :", article)
			reponse.Body.Close()
			continue
		}

		document, err := goquery.NewDocumentFromReader(reponse.Body)
		reponse.Body.Close()
		if err != nil {
			fmt.Println("Erreur HTML :", err)
			continue
		}

		var paragraphes []string

		document.Find("div.mw-parser-output > p").Each(func(index int, element *goquery.Selection) {
			texte := strings.TrimSpace(element.Text())
			if texte != "" {
				paragraphes = append(paragraphes, texte)
			}
		})

		if len(paragraphes) == 0 {
			fmt.Println("Aucun paragraphe trouvé pour", article)
			continue
		}

		totalMots := 0
		for j := 0; j < len(paragraphes); j++ {
			mots := strings.Fields(paragraphes[j])
			totalMots += len(mots)
		}

		fmt.Println("Article :", article)
		fmt.Println("Nb paragraphes :", len(paragraphes))
		fmt.Println("Nb mots :", totalMots)

		os.MkdirAll("out", os.ModePerm)

		nomFichier := "out/wiki_" + article + ".txt"

		var contenu []string
		contenu = append(contenu, "Article : "+article)
		contenu = append(contenu, "Nb paragraphes : "+strconv.Itoa(len(paragraphes)))
		contenu = append(contenu, "Nb mots : "+strconv.Itoa(totalMots))
		contenu = append(contenu, "Contenu : ")

		for j := 0; j < len(paragraphes); j++ {
			contenu = append(contenu, paragraphes[j])
		}

		ecrireFichier(nomFichier, contenu)

		fmt.Println("Fichier généré :", nomFichier)
	}
}
