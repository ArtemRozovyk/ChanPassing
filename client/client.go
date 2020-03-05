package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	st "./structures" // contient la structure Personne
	tr "./travaux"
)

var ADRESSE string = "localhost"                           // adresse de base pour la Partie 2
var FICHIER_SOURCE string = "./conseillers-municipaux.txt" // fichier dans lequel piocher des personnes
var TAILLE_SOURCE int = 450000                             // inferieure au nombre de lignes du fichier, pour prendre une ligne au hasard
var TAILLE_G int = 5                                       // taille du tampon des gestionnaires
var NB_G int = 2                                           // nombre de gestionnaires
var NB_P int = 2                                           // nombre de producteurs
var NB_O int = 4                                           // nombre d'ouvriers
var NB_PD int = 2                                          // nombre de producteurs distants pour la Partie 2
var lnnb int = 2
var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "M"} // une personne vide

// paquet de personne, sur lequel on peut travailler, implemente l'interface personne_int

type personne_emp struct {
	pers    st.Personne
	ligne   int
	lecteur chan message_lec
	afaire  []func(st.Personne) st.Personne
	status  string
}

type message_lec struct {
	retour chan string
	ligne  int
}

// paquet de personne distante, pour la Partie 2, implemente l'interface personne_int
type personne_dist struct {
	// A FAIRE
}

// interface des personnes manipulees par les ouvriers, les
type personne_int interface {
	initialise()          // appelle sur une personne vide de statut V, remplit les champs de la personne et passe son statut à R
	travaille()           // appelle sur une personne de statut R, travaille une fois sur la personne et passe son statut à C s'il n'y a plus de travail a faire
	vers_string() string  // convertit la personne en string
	donne_statut() string // renvoie V, R ou C
}

// fabrique une personne à partir d'une ligne du fichier des conseillers municipaux
// à changer si un autre fichier est utilisé
func personne_de_ligne(l string) st.Personne {
	//print(l+"\n")
	var separateur = regexp.MustCompile("\u0009") // oui, les donnees sont separees par des tabulations ... merci la Republique Francaise
	separation := separateur.Split(l, -1)
	naiss, _ := time.Parse("2/1/2006", separation[7])
	a1, _, _ := time.Now().Date()
	a2, _, _ := naiss.Date()
	agec := a1 - a2
	return st.Personne{Nom: separation[4], Prenom: separation[5], Sexe: separation[6], Age: agec}
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES ***

func (p *personne_emp) initialise() {
	// A FAIRE
	recupLine := make(chan string)
	p.lecteur <- message_lec{retour: recupLine, ligne: p.ligne}
	line := <-recupLine
	p.pers = personne_de_ligne(line)
	p.status = "R"
	for i := 0; i < 5; i++ {
		p.afaire = append(p.afaire, tr.UnTravail())
	}

}

func (p *personne_emp) travaille() {
	var tache func(st.Personne) st.Personne
	tache, p.afaire = p.afaire[0], p.afaire[1:]
	tache(p.pers)
	if len(p.afaire) == 0 {
		p.status = "C"
	} else {
		print(len(p.afaire))
		//fmt.Println("Length : ",len(p.afaire), p.pers.Nom,p.pers.Prenom)

	}
}

func (p *personne_emp) vers_string() string {
	return p.pers.Nom + " " + p.pers.Prenom + " " + p.status
	//TODO
}

func (p *personne_emp) donne_statut() string {
	return p.status
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES DISTANTES (PARTIE 2) ***
// ces méthodes doivent appeler le proxy (aucun calcul direct)

func (p personne_dist) initialise() {
	// A FAIRE
}

func (p personne_dist) travaille() {
	// A FAIRE
}

func (p personne_dist) vers_string() string {
	// A FAIRE
	return ""
}

func (p personne_dist) donne_statut() string {
	// A FAIRE
	return ""
}

// *** CODE DES GOROUTINES DU SYSTEME ***

// Partie 2: contacté par les méthodes de personne_dist, le proxy appelle la méthode à travers le réseau et récupère le résultat
// il doit utiliser une connection TCP sur le port donné en ligne de commande
func proxy() {
	// A FAIRE
}

// Partie 1 : contacté par la méthode initialise() de personne_emp, récupère une ligne donnée dans le fichier source
func lecteur(recupLine chan message_lec) {
	for {
		lNb := <-recupLine
		fmt.Println("Lecteur solicité pour la ligne ", lNb.ligne)

		file, err := os.Open("./client/1-rne-cm.txt")
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(file)
		for i := 0; i < lNb.ligne; i++ {
			_ = scanner.Scan()
		}

		line := scanner.Scan()
		if line == false {
			log.Fatal(err)
		} else {
			r := scanner.Text()
			fmt.Print("La ligne trouvé:  ", r)
			lNb.retour <- r
		}

		file.Close()
	}

}

// Partie 1: récupèrent des personne_int depuis les gestionnaires, font une opération dépendant de donne_statut()
// Si le statut est V, ils initialise le paquet de personne puis le repasse aux gestionnaires
// Si le statut est R, ils travaille une fois sur le paquet puis le repasse aux gestionnaires
// Si le statut est C, ils passent le paquet au collecteur

// Partie 1: les producteurs cree des personne_int implementees par des personne_emp initialement vides,
// de statut V mais contenant un numéro de ligne (pour etre initialisee depuis le fichier texte)
// la personne est passée aux gestionnaires
func producteur(gestChan chan personne_int, lec chan message_lec) {
	var b = 1
	for {
		initPack := personne_emp{pers_vide, b, lec, make([]func(st.Personne) st.Personne, 0), "V"}
		b++
		gestChan <- &initPack
	}

	// A FAIRE
}

// Partie 2: les producteurs distants cree des personne_int implementees par des personne_dist qui contiennent un identifiant unique
// utilisé pour retrouver l'object sur le serveur
// la creation sur le client d'une personne_dist doit declencher la creation sur le serveur d'une "vraie" personne, initialement vide, de statut V
func producteur_distant() {
	// A FAIRE

}

// Partie 1: les gestionnaires recoivent des personne_int des producteurs et des ouvriers et maintiennent chacun une file de personne_int
// ils les passent aux ouvriers quand ils sont disponibles
// ATTENTION: la famine des ouvriers doit être évitée: si les producteurs inondent les gestionnaires de paquets, les ouvrier ne pourront
// plus rendre les paquets surlesquels ils travaillent pour en prendre des autres

func gestionnaire(fromProducer chan personne_int, toOuvr chan personne_int, fromOuv chan personne_int, ack chan int) {
	//var queue [] personne_int
	queue := make([]personne_int, 0)
	for {
		switch len(queue) {
		case TAILLE_G:
			toOuvr <- queue[0]
			queue = queue[1:]
		case 0:
			select {
			case np := <-fromOuv:
				queue = append(queue, np)
			case np := <-fromProducer:
				queue = append(queue, np)
			}

		default:
			if len(queue) < TAILLE_G/2 {
				select {
				case np := <-fromOuv:
					queue = append(queue, np)
				case np := <-fromProducer:
					queue = append(queue, np)
				case toOuvr <- queue[0]:
					queue = queue[1:]
				}
			} else {
				select {
				case np := <-fromOuv:
					queue = append(queue, np)
				case toOuvr <- queue[0]:
					queue = queue[1:]
				}
			}
		}
	}
}

func ouvrier(fromGest chan personne_int, togest chan personne_int, collect chan personne_int, ack chan int) {
	for {
		pers := <-fromGest
		switch pers.donne_statut() {
		case "V":
			println("got V from gest")
			pers.initialise()
			println("sending to gest " + pers.vers_string())
			togest <- pers
			//ack<-0
			println("sent to gest " + pers.vers_string())
		case "R":
			println("got R from gest")
			pers.travaille()
			togest <- pers
			println("sent R to gest " + pers.vers_string())
		case "C":
			println("sent to collect")
			collect <- pers
		}
	}
	// A FAIRE
}

// Partie 1: le collecteur recoit des personne_int dont le statut est c, il les collecte dans un journal
// quand il recoit un signal de fin du temps, il imprime son journal.
func collecteur(fromOuvr chan personne_int) {
	for {
		println("DONE " + (<-fromOuvr).vers_string())
	}
	// A FAIRE
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // graine pour l'aleatoire
	if len(os.Args) < 3 {
		fmt.Println("Format: client <port> <millisecondes d'attente>")
		return
	}
	//port, _ := strconv.Atoi(os.Args[1]) // utile pour la partie 2
	millis, _ := strconv.Atoi(os.Args[2]) // duree du timeout
	fintemps := make(chan int)
	prodGest := make(chan personne_int)
	gestOuvr := make(chan personne_int)
	gestOuvr2 := make(chan personne_int)
	ouvrCollect := make(chan personne_int)
	toLec := make(chan message_lec)
	ack := make(chan int)
	go func() { producteur(prodGest, toLec) }()
	go func() { gestionnaire(prodGest, gestOuvr, gestOuvr2, ack) }()
	go func() { ouvrier(gestOuvr, gestOuvr2, ouvrCollect, ack) }()
	//go func() {producteur(ouvrCollect)}()
	go func() { collecteur(ouvrCollect) }()
	go func() { lecteur(toLec) }()
	// A FAIRE
	// creer les canaux
	// lancer les goroutines (parties 1 et 2): 1 lecteur, 1 collecteur, des producteurs, des gestionnaires, des ouvriers
	// lancer les goroutines (partie 2): des producteurs distants, un proxy
	time.Sleep(time.Duration(millis) * time.Millisecond)
	fintemps <- 0
	<-fintemps
}
