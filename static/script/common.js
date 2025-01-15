function declareTagInputForElements(tagsContainer, textInput, hiddenInput) {
    const tags = document.getElementById(tagsContainer);
    const input = document.getElementById(textInput);

    input.addEventListener('keydown', function (event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            const tag = document.createElement('li');
            const tagContent = input.value.trim();
            if (tagContent !== '') {
                tag.innerText = tagContent;
                tag.innerHTML += '<button class="delete-button">X</button>';
                tags.appendChild(tag);
                input.value = '';
                document.querySelectorAll('input[name*="'+ hiddenInput +'"]').forEach((element) => {
                    element.remove();
                })
                Array.from(document.querySelectorAll("#" + tagsContainer + " > li")).forEach((value, index) => {
                    const hiddenElement = Object.assign(document.createElement('input'),
                        {type: 'hidden', value: value.firstChild.textContent, name: hiddenInput + '[' + index + ']'});
                    tags.appendChild(hiddenElement);
                })
            }
        }
    });

    tags.addEventListener('click', function (event) {
        if (event.target.classList.contains('delete-button')) {
            event.target.parentNode.remove();
            document.querySelectorAll('input[name*="'+ hiddenInput +'"]').forEach((element) => {
                element.remove();
            })
            var tagValues = Array.from(document.querySelectorAll("#" + tagsContainer + " > li")).forEach((value, index) => {
                const hiddenElement = Object.assign(document.createElement('input'),
                    {type: 'hidden', value: value.firstChild.textContent, name: hiddenInput + '[' + index + ']'});
                tags.appendChild(hiddenElement);
            })
        }
    });
}

function declareMarkdownPreview(sourceInput, targetElement) {
    const source = document.getElementById(sourceInput);
    const target = document.getElementById(targetElement);

    source.addEventListener('input', function () {
       target.innerHTML = marked.parse(source.value);
    })

    target.innerHTML = marked.parse(source.value);
}

function showActiveTreeSelection(expandableClassName, selectedClassName) {
    const selectedElements = document.getElementsByClassName(selectedClassName);
    const urlParams = new URLSearchParams(window.location.search);
    const treePath = urlParams.get("path");

    console.log(selectedElements);

    Array.from(selectedElements).forEach((element) => {
        let breadcrumbs = element.attributes["breadcrumbs"];
        if (breadcrumbs && breadcrumbs.value == treePath) {
            element.classList.add("selected-tree-node");
        }
    });

    const expandableElements = document.getElementsByClassName(expandableClassName);
    Array.from(expandableElements).forEach((element) => {
        let breadcrumbs = element.attributes["breadcrumbs"];
        if (breadcrumbs && treePath.startsWith(breadcrumbs.value)) {
            element.open = true;
        }
    });
}