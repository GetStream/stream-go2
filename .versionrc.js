const versionFileUpdater = {
    REGEX: /Version = "v([0-9.]+)"/,

    readVersion: function (contents) {
        return this.REGEX.exec(contents)[1];
    },

    writeVersion: function (contents, version) {
        const splitted = version.split('.');
        const [major, minor, patch] = [splitted[0], splitted[1], splitted[2]];

        return contents.replace(this.REGEX, `Version = "v${major}.${minor}.${patch}"`);
    }
}

const moduleVersionUpdater = {
    GO_MOD_REGEX: /stream-go2\/v(\d+)/g,

    readVersion: function (contents) {
        return this.GO_MOD_REGEX.exec(contents)[1];
    },

    writeVersion: function (contents, version) {
        const major = version.split('.')[0];
        const previousMajor = major - 1;
        const go_mod_regex = new RegExp(`stream-go2\/v(${previousMajor})`, "g");

        return contents.replace(go_mod_regex, `stream-go2/v${major}`);
    }
}

module.exports = {
    bumpFiles: [
        { filename: './version.go', updater: versionFileUpdater },
        { filename: './go.mod', updater: moduleVersionUpdater },
        { filename: './README.md', updater: moduleVersionUpdater },
    ],
}