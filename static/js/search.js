oc = oc || {};

oc.search = oc.search || {};

oc.search.page = function(text) {
    $.ajax({
        url: "/search/" + text
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data.Error);
        }
        if (data.Result) {
            window.location.href = data.Result;
        }
    });
};