oc = oc || {};

oc.nav = oc.nav || {};

var page = "";

$(document).ready(function() {
    page = $('#data-store').data('page')
});

$(document).on('click', 'a', function (e) {
    e.preventDefault();
    oc.nav.move($(this)[0].pathname);
});

oc.nav.move = function(location) {
    if (page !== "") {
        if (location.lastIndexOf("/wiki/", 0) !== 0) {
            if (location.lastIndexOf("/user/", 0) === 0) {
                window.location.href = location;
            }
            return;
        }
        oc.chat.sendMessage("[[[" + location + "]]]");
    }
    window.location.href = location;
};