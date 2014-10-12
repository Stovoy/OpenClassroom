oc = oc || {};

oc.search = oc.search || {};

oc.search.page = function(text, sel, redirect, callback) {
    text = text.replace(/ /g, "+");
    $.ajax({
        url: "/search/" + text
    }).done(function(data) {
        data = JSON.parse(data);
        var $searchInput = $(sel);
        $searchInput.css('background-color', 'white');
        if (data.Error) {
            $searchInput.css('background-color', 'red');
        }
        if (data.Result) {
            if (redirect) {
                oc.nav.move(data.Result);
            }
        }
        if (callback) {
            callback(data);
        }
    });
};