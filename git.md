# Git Remotes

These repositories use GitHub as the main remote for pulling/fetching, while pushes are sent to both GitHub and CTI GitLab.

After cloning either repository, run the matching commands below from inside the cloned repository.

## maas-ldap

```bash
git remote set-url origin git@github.com:TudorBogos/maas-ldap.git
git remote set-url --add --push origin git@github.com:TudorBogos/maas-ldap.git
git remote set-url --add --push origin git@git.cti.ugal.ro:tb171/maas-ldap.git
```

Expected result:

```text
origin  git@github.com:TudorBogos/maas-ldap.git (fetch)
origin  git@github.com:TudorBogos/maas-ldap.git (push)
origin  git@git.cti.ugal.ro:tb171/maas-ldap.git (push)
```

## Check

```bash
git remote -v
git remote get-url origin
git remote get-url --push --all origin
```

## Usage

Pull from GitHub only:

```bash
git pull
```

Push to both GitHub and CTI GitLab:

```bash
git push
```

Note: pushing to multiple remotes is not atomic. If one remote succeeds and the other fails, fix the failing remote and run `git push` again.
