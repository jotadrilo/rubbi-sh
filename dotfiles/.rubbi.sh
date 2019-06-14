# shellcheck disable=SC2148

alias rubsh="rubbi-sh"
alias rubclean="rubsh -clean"
alias rubshow="rubsh -show"
alias rubhelp="rubsh -help"
alias rubbish="rubcd"

function rubdel {
    rubsh -del "${1}"
}
function rubadd {
    rubsh -add "${1}"
}
function rubuse {
    rubsh -use "${1}"
}
function rubsel {
    rubsh -show
    echo
    echo -n "Folder to use: "
    read -r fn
    rubsh -use "$fn"
    rubcd
}
function rubcd {
    cd "$(rubsh)" || exit 1
}
