async function getWikipediaImage(wikiUrl: string): Promise<string | null> {
    const pageTitle = wikiUrl.split('/wiki/')[1];
    const apiUrl = `https://en.wikipedia.org/w/api.php?action=parse&page=${pageTitle}&format=json&prop=text&origin=*`;

    try {
        const response = await fetch(apiUrl);
        const data = await response.json();
        const pageContent = data.parse.text['*'];
        const parser = new DOMParser();
        const doc = parser.parseFromString(pageContent, 'text/html');
        const infobox = doc.querySelector('.infobox');
        const img = infobox?.querySelector('img');
        return img ? `https:${img.getAttribute('src')}` : null;
    } catch (error) {
        console.error('Error fetching Wikipedia image:', error);
        return null;
    }
}

export default getWikipediaImage;