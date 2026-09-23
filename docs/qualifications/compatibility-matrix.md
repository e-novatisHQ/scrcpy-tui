# Matrice de compatibilité — QUAL-01

État au 23 septembre 2026 : **plage de versions non établie**. Les résultats
historiques ci-dessous sont des combinaisons ponctuelles, pas une garantie pour
les versions intermédiaires, d'autres firmwares ou un nouvel hôte.

| Lanceur / hôte | adb / scrcpy | Cible / transport | Preuve | Limite |
|---|---|---|---|---|
| 0.3.1 / Debian 13 amd64 | 1.0.41, platform-tools 34.0.5 / 4.1 | Homatics Android 14 / Wi-Fi | [Compte rendu](homatics-dongle-g-4k-v0.3.1.md) | USB, audio et contrôle non qualifiés |
| Même combinaison | Même combinaison | Strong Android 12 / Wi-Fi | [Compte rendu](strong-uhd-google-tv-stick-v0.3.1.md) | Audio encodé, qualité audible non vérifiée |
| Même combinaison | Même combinaison | Pixel 6 Android 16 / USB ; Xiaomi TV Android 11 / Wi-Fi | [Compte rendu](pixel6-xiaomi-mitv-v0.3.1.md) | Contrôle et longue durée non vérifiés |
| Même combinaison | Même combinaison | Xiaomi Mi A3 Android 11 / USB | [Compte rendu](xiaomi-mi-a3-v0.3.1.md) | **Échec audio** ; vidéo et presets validés |
| CORE-01 / Linux et Windows CI | Helpers synthétiques | Aucun appareil | Tests `internal/app/capabilities_test.go` | Prouve le diagnostic, pas la compatibilité réelle |
| Candidat Windows 0.5.0 | À relever | À identifier | [Protocole non exécuté](windows-11-v0.5.0-protocol.md) | Aucune qualification Windows 11 réelle |

## Protocole de promotion d'une combinaison

1. Fixer le commit du lanceur, son artefact et son SHA-256 ; vérifier les sommes
   publiées et l'attestation avant l'installation. Relever OS, architecture,
   versions complètes et chemins d'adb/scrcpy, source et digest des outils.
2. Enregistrer l'aide `scrcpy --help` avec sa version et son digest. Exécuter les
   tests de capacités ; conserver le résultat distinct des essais matériels.
3. Identifier un appareil autorisé, modèle, firmware, version Android et transport.
   Employer un profil temporaire, sans remplacer celui de l'opérateur. Ne publier
   ni numéro de série, adresse ADB, capture d'écran privée ou enregistrement audio.
4. Exécuter les trois presets, un argument cité, disparition avant lancement,
   interruption, retour au terminal et contrôle des seuls descendants possédés.
   Vérifier visuellement la vidéo. Tester USB et Wi-Fi séparément ; noter NOT RUN
   si l'un n'est pas disponible. Tester l'audio et le contrôle uniquement dans le
   périmètre autorisé, et consigner leur résultat séparément.
5. Répéter après redémarrage de l'hôte et de l'appareil. Conserver pour chaque essai
   heure, commande expurgée, résultat attendu/observé, code de sortie, preuve et
   restitution du profil/de la connexion. Un code zéro seul ne prouve pas l'audio.
6. Faire relire les preuves. Promouvoir seulement les combinaisons démontrées.
   Toute plage proposée exige au minimum les bornes et les changements de syntaxe
   pertinents ; ne pas déduire une compatibilité continue d'un seul point testé.

Le choix des versions candidates dépend des versions réellement utilisées par
les opérateurs. Aucun minimum universel ne sera ajouté aux paquets à partir des
seuls helpers. Refaire les essais lorsqu'un preset, le backend de session ou une
version d'adb/scrcpy change. Les logs bruts restent privés ; le compte rendu durable
contient les observations expurgées et les digests des artefacts publics.
