package response

type Response struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
	PriceChange string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	HighPrice string `json:"highPrice"`
	LowPrice string `json:"lowPrice"`
	Volume string `json:"volume"`
}

type Currency struct {
	Symbol string `json:"symbol"`
	Price float64 `json:"price"`
	PriceChange float64 `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	HighPrice float64 `json:"highPrice"`
	LowPrice float64 `json:"lowPrice"`
	Volume float64 `json:"volume"`
}

var DefaultList = "btc,eth,bnb,bcc,neo,ltc,qtum," +
	"ada,xrp,eos,tusd,iota,xlm,ont,trx,etc,icx,ven," +
	"nuls,vet,pax,bchabc,bchsv,usdc,link,waves,btt,usds,ong,hot,"+
	"zil,zrx,fet,bat,xmr,zec,iost,celr,dash,nano,omg,theta,enj,mith,"+
	"matic,atom,tfuel,one,ftm,algo,usdsb,gto,erd,doge,dusk,ankr,win,"+
	"cos,npx,cocos,mtl,tomo,perl,dent,mft,key,storm,dock,wan,fun,cvc,"+
	"chz,band,busd,beam,xtz,ren,rvn,hc,hbar,nkn,stx,kava,arpa,iotx,rlc,"+
	"mco,ctxc,bch,troy,vite,ftt,eur,ogn,drep,bull,bear,ethbull,ethbear,tct,"+
	"wrx,bts,lsk,bnt,lto,eosbull,eosbear,xrpbull,xrpbear,strat,aion,mbl,coti,"+
	"bnbbull,bnbbear,stpt,wtc,data,xzc,sol,ctsi,hive,chr,btcup,btcdown,gxs,"+
	"ardr,lend,mdt,stmx,knc,rep,lrc,pnt,comp,bkrw,sc,zen,snx,ethup,ethdown,"+
	"adaup,adadown,linkup,linkdown,vtho,dgb,gbp,sxp,mkr,dai,dcr,storj,bnbup,"+
	"bnbdown,xtzup,xtzdown,mana,aud,yfi,bal,blz,iris,kmd,jst,srm,ant,crv,"+
	"sand,ocean,nmr,dot,luna,rsr,paxg,wnxm,trb,bzrx,sushi,fyi,ksm,egld,dia,"+
	"rune,fio,uma,eosup,eosdown,trxup,trxdown,xrpup,xrpdown,dotup,dotdown,"+
	"bel,wing,ltcup,ltcdown,uni,nbs,oxt,sun,avax,hnt,flm,uniup,unidown,"+
	"orn,utk,xvs,alpha,aave,near,sxpup,sxpdown,fil,filup,fildown,yfiup,yfidown,"+
	"inj,audio,ctk,bchup,bchdown,akro,axs,hard,dnt,strax,unfi,rose,ava,xem,"+
	"aaveup,aavedown,skl,susd,sushiup,sushidown,xlmup,xlmdown,grt,juv,psg,1inch,"+
	"reef,og,atm,asr,celo,rif,btcst,tru,ckb,twt,firo,lit,sfp,dodo,cake,acm,"+
	"badger,fis,om,pond,dego,alice,lina,perp,ramp,super,cfx,eps,auto,tko,"+
	"pundix,tlm,1inchup,1inchdown,btg,mir,bar,forth,bake,burger,slp,shib,icp,"+
	"ar,pols,mdx,mask,lpt,nu,xvg,ata,gtc,torn,keep,ern,klay,pha,bond,mln,dexe,"+
	"c98,clv,qnt,flow,tvk,mina,ray,farm,alpaca,quick,mbox,for,req,ghst,waxp,tribe"