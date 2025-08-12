# Fish shell version of .settings.sh for p2k16

# Set basedir to the directory of this script
set basedir (dirname (status --current-filename))
cd $basedir
set basedir (pwd)
cd - > /dev/null

# Enable NVM if available
test -r $HOME/.nvm/nvm.sh; and begin
    echo "Loading NVM"
    source $HOME/.nvm/nvm.sh
    nvm use
end

# Add bin/ to PATH
set -gx PATH $basedir/bin $PATH

echo "Setting PGPORT and PGHOST."
set -gx PGHOST (or $PGHOST "127.0.0.1")
set -gx PGPORT (or $PGPORT "2016")
set -gx PGPASSFILE (or $PGPASSFILE (pwd)/.pgpass)
if test -f $PGPASSFILE
    chmod 0600 $PGPASSFILE
end

echo "Setting FLYWAY_*"
set -gx FLYWAY_URL jdbc:postgresql://$PGHOST:$PGPORT/p2k16
set -gx FLYWAY_USER p2k16-flyway
set -gx FLYWAY_PASSWORD p2k16-flyway
set -gx FLYWAY_SCHEMAS public
set -gx FLYWAY_VALIDATE_ON_MIGRATE false
