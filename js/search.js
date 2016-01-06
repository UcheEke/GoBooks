/**
 * Created by ekeu on 06/01/16.
 */

window.onload(
    function submitSearch() {
        $.ajax({
            url:"/search",
            method:"POST",
            data:$("#searchForm")
        });
    }
)();