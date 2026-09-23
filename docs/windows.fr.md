# Windows — installation et diagnostic

Guide en préparation pour v0.5.0, cible Windows 10/11 x64. La release publique
actuelle v0.4.0 reste Linux uniquement. Ne considérer un paquet Windows comme
qualifié qu’après validation des gates du [plan](project-status-and-plan.md).
Windows ARM64 et les sessions sans console ne sont pas qualifiés.

## Prérequis

Utiliser PowerShell dans Windows Terminal. Installer séparément
[scrcpy pour Windows](https://github.com/Genymobile/scrcpy/blob/master/doc/windows.md)
et adb, fourni dans les distributions scrcpy officielles ou dans les
[Platform-Tools Android](https://developer.android.com/tools/releases/platform-tools).
Le MSI du lanceur ne les installe pas et ne télécharge rien au lancement.
Garder les DLL et fichiers accompagnant scrcpy dans son répertoire d’origine.
Ajouter les répertoires contenant `scrcpy.exe` et `adb.exe` au PATH utilisateur,
puis fermer et rouvrir Windows Terminal. Vérifier :

```powershell
Get-Command adb.exe, scrcpy.exe | Select-Object Name, Source
adb.exe version
scrcpy.exe --version
```

En cas de plusieurs installations, `Get-Command -All adb.exe, scrcpy.exe` montre
les candidats : conserver un PATH cohérent et relancer le terminal. Ne pas
remplacer un adb utilisé par d’autres outils pendant leurs opérations.

## Installer et mettre à niveau

Quand la release qualifiée est disponible, télécharger le MSI x64 et `SHA256SUMS`
de la **même release**, puis comparer l’empreinte avant d’ouvrir le MSI :

```powershell
Get-FileHash .\scrcpy-tui_0.5.0_windows_amd64.msi -Algorithm SHA256
Select-String -Path .\SHA256SUMS -Pattern 'scrcpy-tui_0.5.0_windows_amd64.msi$'
Get-AuthenticodeSignature .\scrcpy-tui_0.5.0_windows_amd64.msi
```

Le nom ci-dessus est un exemple pour la future v0.5.0, pas une annonce de
publication. Un checksum confirme l’intégrité par rapport au manifeste choisi ;
il ne constitue pas une signature Authenticode. Les paquets de développement
actuels ne sont pas signés et peuvent déclencher SmartScreen. Ne pas désactiver
les protections Windows ; faire valider leur provenance selon la politique locale.
Avec GitHub CLI, une release attestée permet aussi :

```powershell
gh attestation verify .\scrcpy-tui_0.5.0_windows_amd64.msi --repo e-novatisHQ/scrcpy-tui
```

Double-cliquer sur le MSI, depuis le compte qui utilisera l’application. Le paquet
installe le lanceur dans `%LOCALAPPDATA%\Programs\scrcpy-tui` et ajoute ce seul
répertoire au PATH utilisateur. Il ne demande pas de déplacement manuel de fichiers.
Fermer et rouvrir **toutes** les fenêtres Windows Terminal après l’installation,
puis vérifier `scrcpy-tui --version`. Une nouvelle version MSI remplace l’ancienne
installation de ce compte. Fermer les sessions scrcpy et le lanceur avant la mise
à niveau. Les profils utilisateur restent conservés.

Le ZIP constitue une alternative portable : comparer son SHA-256, l’extraire dans
un dossier utilisateur et lancer `scrcpy-tui.exe` depuis PowerShell. Le ZIP ne
modifie pas le PATH. Éviter de mélanger le ZIP et le MSI dans le même répertoire.

## Premier lancement

Activer le débogage Android et autoriser cet ordinateur sur l’appareil. Vérifier
`adb.exe devices -l` : l’état doit être `device`, et non `unauthorized` ou `offline`.
En USB, utiliser un câble de données et le pilote Windows adapté à l’appareil.
En Wi-Fi, établir préalablement la connexion ADB autorisée : le lanceur ne configure
ni le réseau ni l’appairage.

```powershell
scrcpy-tui devices
scrcpy-tui presets
scrcpy-tui preview --device 'USB_TEST' --preset 'Très léger'
scrcpy-tui launch --device 'USB_TEST' --preset 'Très léger' --yes
scrcpy-tui
```

Remplacer `USB_TEST` par la cible affichée localement ; expurger cet identifiant des
preuves partagées. `preview` ne lance rien ; `launch --yes` revalide l’état ADB.
Les commandes et les [touches du guide général](usage.fr.md) restent identiques.
La commande affichée est une prévisualisation des arguments ; utiliser le CLI
ci-dessus pour éviter les différences d’échappement entre Bash et PowerShell.

Pendant une session, fermer la fenêtre scrcpy ou appuyer sur Ctrl+C dans le
terminal. Le lanceur demande un arrêt gracieux ciblé puis force uniquement son
Job Object après deux secondes si nécessaire. Le serveur adb préexistant n’est
pas arrêté. Le menu revient après la session ; `q` quitte la TUI. Les dernières
sorties scrcpy sont disponibles avec `l`, limitées à 32 Kio en mémoire.

## Diagnostic et récupération

| Symptôme | Vérification ou action |
|---|---|
| Commande introuvable | Rouvrir Windows Terminal ; vérifier `Get-Command -All scrcpy-tui.exe` et le PATH utilisateur |
| adb ou scrcpy introuvable | Vérifier leurs `.exe`, fichiers associés et répertoires PATH avec les commandes ci-dessus |
| Appareil unauthorized/offline | Autoriser le PC sur Android et vérifier le transport avec adb avant de rafraîchir avec `r` |
| Option scrcpy inconnue | Consulter `scrcpy.exe --help`, relever sa version et adapter le preset ; aucune plage de versions Windows n’est encore revendiquée |
| Erreur de Job Object | Conserver le message, version Windows et terminal ; aucune exécution sans isolation n’est tentée |
| Ctrl+C nécessite l’arrêt forcé | Vérifier le retour au menu et l’absence de descendants ; signaler le résultat de la qualification du terminal et de scrcpy |
| Écriture du profil refusée | Vérifier le dossier du profil et ses droits ; une configuration illisible n’est pas écrasée |

Le profil par défaut est `%APPDATA%\scrcpy-tui\presets.json`. Pour une recette
isolée : `scrcpy-tui --config "$env:TEMP\scrcpy-tui-recette\presets.json"`.
Le fichier et ses sauvegardes ont une ACL protégée limitée à l’utilisateur courant
et SYSTEM. Un profil invalide est copié dans `presets.json.recovery-*` avant
réécriture. Fermer le lanceur, conserver le fichier courant, puis recopier la
sauvegarde choisie vers `presets.json` pour restaurer. Les arguments ne doivent
contenir aucun secret : ils restent visibles dans la configuration et les processus.

Désinstaller via Paramètres → Applications → Applications installées → scrcpy-tui,
ou avec le MSI d’origine. Le binaire et son entrée PATH sont retirés ; les profils
sont conservés pour une réinstallation. Rouvrir Windows Terminal. Si une autre
copie du ZIP reste accessible, `Get-Command -All scrcpy-tui.exe` permet de la repérer.
Ne pas supprimer adb ou scrcpy pour désinstaller uniquement le lanceur.
