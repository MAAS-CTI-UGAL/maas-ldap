# Git Remotes

CTI is the main source of truth.

- Pull/fetch from CTI: `git@git.cti.ugal.ro:maas/maas-ldap.git`
- Push to CTI and GitHub fork

## Clone

```bash
git clone git@git.cti.ugal.ro:maas/maas-ldap.git
cd maas-ldap

git remote set-url --push origin git@git.cti.ugal.ro:maas/maas-ldap.git
git remote set-url --add --push origin git@github.com:TudorBogos/maas-ldap.git
```

## Existing Clone

```bash
git remote set-url origin git@git.cti.ugal.ro:maas/maas-ldap.git
git remote set-url --push origin git@git.cti.ugal.ro:maas/maas-ldap.git
git remote set-url --add --push origin git@github.com:TudorBogos/maas-ldap.git
```

## Check

```bash
git remote -v
```

Expected:

```text
origin    git@git.cti.ugal.ro:maas/maas-ldap.git (fetch)
origin    git@git.cti.ugal.ro:maas/maas-ldap.git (push)
origin    git@github.com:TudorBogos/maas-ldap.git (push)
```

## Daily Use

```bash
git pull --ff-only
git push origin dev
```

`git pull --ff-only` pulls from CTI. `git push origin dev` pushes to both CTI and GitHub.
