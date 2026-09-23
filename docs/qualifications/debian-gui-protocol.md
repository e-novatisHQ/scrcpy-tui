# Installation graphique Debian/Ubuntu — QUAL-02

**NOT RUN** — aucun hôte graphique propre et autorisé n'est disponible dans cette
session. L'extraction du `.deb`, la simulation APT, l'association Discover et les
smoke tests existants ne valent pas installation graphique authentifiée.

## Préparation

Utiliser une VM jetable avec instantané, ou un poste de qualification explicitement
réservé. Relever distribution, version, architecture, gestionnaire graphique,
version du paquet et SHA-256. Photographier ou capturer seulement le bureau de
qualification, sans information personnelle. Conserver l'état initial des paquets
`adb`, `scrcpy-tui`, la résolution de la commande et un éventuel profil préexistant.
Vérifier checksum et attestation de l'artefact publié avant de l'ouvrir.

## Recette à exécuter

| Étape | Résultat attendu | Preuve à joindre |
|---|---|---|
| Ouvrir le `.deb` dans le gestionnaire graphique | Nom, version, architecture cohérents | Capture de la fiche |
| Installer et authentifier localement | Installation achevée sans erreur ; dépendance adb résolue | Capture finale et journal expurgé |
| Ouvrir un nouveau terminal | `command -v scrcpy-tui` résout `/usr/bin/scrcpy-tui` ; version correcte | Sorties et version du shell |
| Lancer sans scrcpy disponible dans PATH | Diagnostic utile, terminal restauré | Commande et observation |
| Installer/configurer un scrcpy choisi pour la qualification | Prérequis identifié, version relevée | Source et digest si disponibles |
| Tester un lancement sur appareil autorisé | Résultat consigné séparément par transport | Lien vers compte rendu matériel |
| Désinstaller depuis le gestionnaire graphique | Exécutable du paquet retiré ; aucune suppression implicite du profil utilisateur | Capture et inspection du contenu du paquet |
| Réinstaller puis mettre à niveau vers un candidat plus récent | Version cible correcte, profil utilisateur conservé | État avant/après et logs |
| Restaurer l'instantané | Poste dans son état initial | Heure et confirmation opérateur |

L'ordre de mise à niveau doit utiliser deux versions réelles identifiées ; ne pas
inventer une version publiée pour satisfaire la recette. Le mot de passe reste
saisi par l'opérateur dans le dialogue système, jamais dans les logs ou Plane.
Ne pas supprimer automatiquement des dépendances partagées. Un launcher présent
ailleurs dans PATH doit être distingué du fichier du paquet lors de la désinstallation.

## Décision

Un dossier par distribution et architecture réellement testée, avec attendu,
observé et PASS/FAIL/NOT RUN pour chaque ligne. Un échec d'installation, d'intégrité,
de mise à niveau ou de désinstallation bloque la qualification de cette combinaison.
Les essais réussis sur Debian ne qualifient pas Ubuntu, et inversement. La preuve
matérielle est indépendante de la validation du gestionnaire de paquets.
