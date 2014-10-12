oc = oc || {};

oc.search = oc.search || {};

oc.search.page = function(text) {
    $.ajax({
        url: "/search/" + text
    }).done(function(data) {
        data = JSON.parse(data);
        var $searchInput = $('#search-text');
        $searchInput.css('background-color', 'white');
        if (data.Error) {
            $searchInput.css('background-color', 'red');
        }
        if (data.Result) {
            window.location.href = data.Result;
        }
    });
};