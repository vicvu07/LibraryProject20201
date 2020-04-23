var doubleClickPrevent = (dom, fn, e) => {
    dom.prop("disabled", true);

    fn(e);

    setTimeout(function () {
        dom.prop("disabled", false);
    }, 300);
}

createVueDropBox = target => {
    return new Vue({
        el: target,
        data: {
            selected: 0,
            options: []
        }
    });
}

var initAfterLoad = (dom, state) => {
    //  Activate the Tooltips
    dom.find('[data-toggle="tooltip"], [rel="tooltip"]').tooltip();

    // Activate Popovers and set color for popovers
    dom.find('[data-toggle="popover"]').each(function () {
        color_class = $(this).data('color');
        $(this).popover({
            template: '<div class="popover ' + color_class + ' " role="tooltip"><h3 class="popover-title"></h3><div class="popover-content"></div></div>'
        });
    });

    $navbar = dom.find('.navbar[color-on-scroll]');
    scroll_distance = $navbar.attr('color-on-scroll') || 500;

    // Check if we have the class "navbar-color-on-scroll" then add the function to remove the class "navbar-transparent" so it will transform to a plain color.

    if (dom.find('.navbar[color-on-scroll]').length != 0) {
        nowuiKit.checkScrollForTransparentNavbar();
        $(window).on('scroll', nowuiKit.checkScrollForTransparentNavbar)
    }

    dom.find('.form-control').on("focus", function () {
        $(this).parent('.input-group').addClass("input-group-focus");
    }).on("blur", function () {
        $(this).parent(".input-group").removeClass("input-group-focus");
    });

    // Activate bootstrapSwitch
    dom.find('.bootstrap-switch').each(function () {
        $this = $(this);
        data_on_label = $this.data('on-label') || '';
        data_off_label = $this.data('off-label') || '';

        $this.bootstrapSwitch({
            onText: data_on_label,
            offText: data_off_label
        });
    });

    // Activate Carousel
    dom.find('.carousel').carousel({
        interval: 4000
    });

    dom.find('.date-picker').each(function () {
        $(this).datepicker({
            templates: {
                leftArrow: '<i class="now-ui-icons arrows-1_minimal-left"></i>',
                rightArrow: '<i class="now-ui-icons arrows-1_minimal-right"></i>'
            },
            language: 'vi'
        }).on('show', function () {
            $('.datepicker').addClass('open');

            datepicker_color = $(this).data('datepicker-color');
            if (datepicker_color.length != 0) {
                $('.datepicker').addClass('datepicker-' + datepicker_color + '');
            }
        }).on('hide', function () {
            $('.datepicker').removeClass('open');
        });
    });
}

(function ($) {
    $.fn.serializeFormJSON = function () {
        var o = {};
        var a = this.serializeArray();
        $.each(a, function () {
            if (o[this.name]) {
                if (!o[this.name].push) {
                    o[this.name] = [o[this.name]];
                }
                o[this.name].push(this.value || '');
            } else {
                o[this.name] = this.value || '';
            }
        });
        return o;
    };
})(jQuery);