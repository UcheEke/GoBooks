
function showViewPage(){
    $("#search-page").hide();
    $("#search-form").hide();
    $("#view-page").show();
}

function showSearchPage(){
    $("#view-page").hide();
    $("#search-page").show();
    $("#search-form").show();
}

function deleteBook(pk){
    $.ajax({
        url: "/books/"+pk,
        method: "DELETE",
        success: function(){
            $("#book-row-" + pk).remove()
        }
    })
}

function submitSearch(){
        $.ajax({
            url: "/search",
            method: "POST",
            data: $("#search-form").serialize(),
            success: function(rawData){
                var parsed = JSON.parse(rawData);
                if (!parsed) return;

                var searchResults = $("#search-results");
                searchResults.empty();  // Remove past results from display

                // Create the new table entries for each entry
                parsed.forEach(function (result) {
                    var row = $("<tr class='table-hover'><td>" + result.Title
                        + "</td><td>" + result.Author
                        + "</td><td>" + result.Year
                        + "</td><td>" + result.ID
                        + "</td></tr>");
                    // Append to the table
                    searchResults.append(row);

                    // Add functionality to each row for click events
                    row.on("click", function(){
                        $.ajax({
                            url:"/books?id=" + result.ID,
                            method: "PUT",
                            success: function(data) {
                                var book = JSON.parse(data);
                                if (!book) {
                                    return false;
                                }
                                $("#view-results").append("<tr id='book-row-"+ book.PK +"'><td>"+
                                    book.Title+"</td><td>"+
                                    book.Author+"</td><td>" +
                                    book.Classification +
                                    "</td><td><button class='btn btn-danger btn-xs delete-btn' onclick='deleteBook(" + book.PK +")'>" +
                                    "<span class='glyphicon glyphicon-remove-sign'></span>"+
                                    "</button></td></tr>");
                            }
                        })
                    });
                });
            }
        });

    return false;
    }
