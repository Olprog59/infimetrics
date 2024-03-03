(function(){
    function maybeRemoveMe(elt) {
        const timing = elt.getAttribute("remove-content") || elt.getAttribute("data-remove-content");
        if (timing) {
            setTimeout(function () {
                elt.innerHTML = "";
            }, htmx.parseInterval(timing));
        }
    }

    htmx.defineExtension('remove-content', {
        onEvent: function (name, evt) {
            if (name === "htmx:afterProcessNode") {
                var elt = evt.detail.elt;
                if (elt.getAttribute) {
                    maybeRemoveMe(elt);
                    if (elt.querySelectorAll) {
                        var children = elt.querySelectorAll("[remove-content], [data-remove-content]");
                        for (var i = 0; i < children.length; i++) {
                            maybeRemoveMe(children[i]);
                        }
                    }
                }
            }
        }
    });
})();
