#!/bin/bash

# ==========================================
# PortfolioMenu - GitHub Update Script
# ==========================================

set -e


echo ""
echo "PortfolioMenu Git Update"
echo "========================"
echo ""


# Controllo repository

if [ ! -d ".git" ]; then

    echo "Errore: questa cartella non è un repository Git"

    exit 1

fi



# Mostra stato

echo "Stato repository:"
echo "-----------------"

git status

echo ""



# Chiede messaggio commit

read -p "Messaggio commit: " MESSAGE



if [ -z "$MESSAGE" ]; then

    MESSAGE="Update PortfolioMenu"

fi



echo ""

echo "Aggiunta file..."

git add .



echo ""

echo "Commit: $MESSAGE"

git commit -m "$MESSAGE" || echo "Nessun cambiamento da committare"



echo ""

echo "Push GitHub..."

git push origin main



echo ""

echo "================================="
echo "Aggiornamento completato"
echo "================================="
echo ""