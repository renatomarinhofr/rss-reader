export type FeedSource = {
  name: string;
  url: string;
  description?: string;
};

export type FeedCategory = {
  id: string;
  name: string;
  description: string;
  feeds: FeedSource[];
};

export const FEED_CATEGORIES: FeedCategory[] = [
  {
    id: 'noticias',
    name: 'Notícias',
    description: 'Cobertura nacional e regional das principais redações brasileiras.',
    feeds: [
      {
        name: 'G1 - Notícias',
        url: 'https://g1.globo.com/rss/g1/noticias/',
        description: 'Cobertura geral das editorias nacionais do G1.',
      },
      {
        name: 'Agência Brasil',
        url: 'https://agenciabrasil.ebc.com.br/rss/geral/feed.xml',
        description: 'Noticiário oficial da Empresa Brasil de Comunicação (Geral).',
      },
      {
        name: 'BBC News Brasil',
        url: 'https://www.bbc.com/portuguese/index.xml',
        description: 'Principais reportagens internacionais em português.',
      },
      {
        name: 'UOL Notícias',
        url: 'https://rss.uol.com.br/feed/noticias.xml',
        description: 'Manchetes do portal UOL com foco em política e cotidiano.',
      },
      {
        name: 'Canaltech',
        url: 'https://canaltech.com.br/rss/',
        description: 'Tecnologia e inovação com viés jornalístico.',
      },
    ],
  },
  {
    id: 'tecnologia',
    name: 'Tecnologia',
    description: 'Lançamentos, inovação e cultura digital made in Brazil.',
    feeds: [
      {
        name: 'Canaltech',
        url: 'https://canaltech.com.br/rss/',
        description: 'Tecnologia, inovação e análise de produtos.',
      },
      {
        name: 'Tecnoblog',
        url: 'https://tecnoblog.net/feed/',
        description: 'Economia digital, gadgets e bastidores da tecnologia.',
      },
      {
        name: 'G1 - Tecnologia',
        url: 'https://g1.globo.com/rss/g1/tecnologia/',
        description: 'Tendências tecnológicas e ciência aplicada ao dia a dia.',
      },
    ],
  },
  {
    id: 'economia',
    name: 'Economia & Negócios',
    description: 'Mercado, finanças e empreendedorismo com foco local.',
    feeds: [
      {
        name: 'Agência Brasil - Economia',
        url: 'https://agenciabrasil.ebc.com.br/rss/economia/feed.xml',
        description: 'Indicadores econômicos e políticas públicas sob visão oficial.',
      },
      {
        name: 'Valor Econômico',
        url: 'https://valor.globo.com/rss/',
        description: 'Cobertura especializada em finanças, mercado e empresas.',
      },
      {
        name: 'ASN - Agência Sebrae de Notícias',
        url: 'https://agenciasebrae.com.br/feed/',
        description: 'Conteúdo para pequenos negócios e empreendedorismo.',
      },
      {
        name: 'Investing.com Brasil',
        url: 'https://br.investing.com/rss/news.rss',
        description: 'Mercados globais com foco para investidores brasileiros.',
      },
      {
        name: 'InfoMoney',
        url: 'https://www.infomoney.com.br/feed/',
        description: 'Mercado financeiro, investimentos e economia.',
      },
      {
        name: 'Jornal Contábil',
        url: 'https://www.jornalcontabil.com.br/feed/',
        description: 'Tributação, contabilidade e finanças corporativas.',
      },
    ],
  },
  {
    id: 'cultura',
    name: 'Cultura & Entretenimento',
    description: 'Cinema, música, streaming e universo pop sob a ótica nacional.',
    feeds: [
      {
        name: 'G1 - Pop & Arte',
        url: 'https://g1.globo.com/rss/g1/pop-arte/',
        description: 'Cobertura cultural, críticas e agenda do entretenimento.',
      },
      {
        name: 'Omelete',
        url: 'https://www.omelete.com.br/rss',
        description: 'Cinema, séries, HQs e cultura pop.',
      },
      {
        name: 'Chippu',
        url: 'https://chippu.com.br/feed',
        description: 'Curadoria de filmes, séries e cultura pop brasileira.',
      },
      {
        name: 'Revista Arco - Quadrinhos',
        url: 'https://www.ufsm.br/midias/arco/busca?q=&sites%5B0%5D=601&tags=quadrinhos-pt&rss=true',
        description: 'Produções e reportagens sobre HQs e arte gráfica.',
      },
      {
        name: 'Revista Arco - Cultura',
        url: 'https://www.ufsm.br/midias/arco/busca?q=&sites%5B0%5D=601&tags=cultura-pt&rss=true',
        description: 'Reportagens culturais e projetos da UFSM.',
      },
      {
        name: 'SESC SP',
        url: 'https://www.sescsp.org.br/feed/',
        description: 'Agenda cultural, artes e educação do SESC São Paulo.',
      },
      {
        name: 'Revista Continente',
        url: 'http://revistacontinente.com.br/feed',
        description: 'Arte, cultura e sociedade com olhar nordestino.',
      },
      {
        name: 'Mundo Conectado',
        url: 'https://mundoconectado.com.br/feed',
        description: 'Entretenimento, tecnologia e cultura digital.',
      },
    ],
  },
  {
    id: 'ciencia',
    name: 'Ciência & Curiosidades',
    description: 'Descobertas, saúde e histórias que despertam o lado curioso.',
    feeds: [
      {
        name: 'G1 - Ciência e Saúde',
        url: 'https://g1.globo.com/rss/g1/ciencia-e-saude/',
        description: 'Saúde pública, avanços científicos e pesquisas nacionais.',
      },
      {
        name: 'Superinteressante',
        url: 'https://super.abril.com.br/feed/',
        description: 'Ciência, comportamento e curiosidades.',
      },
      {
        name: 'Mega Curioso',
        url: 'https://www.megacurioso.com.br/rss',
        description: 'Histórias incríveis e fatos científicos em linguagem acessível.',
      },
    ],
  },
];
