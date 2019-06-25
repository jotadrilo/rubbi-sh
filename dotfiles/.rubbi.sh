# shellcheck disable=SC2148 disable=SC1090

alias rubsh='rubbi-sh'
alias rubclean='rubbi-sh -clean'
alias rubshow='rubbi-sh -show'
alias rubhelp='rubbi-sh -help'

# Use the handiest alias you prefer
alias r='rbsh'
alias rubbish='rbsh'
alias rubcd='rbsh'

function rubdel {
    rubbi-sh -del
}
function rubadd {
    rubbi-sh -add "${1}"
}
function rubsel {
    rubbi-sh -sel
    rubcd
}
function rbsh {
    cd "$(rubbi-sh)" || exit 1
}
