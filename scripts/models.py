class EndpointTest:
    def __init__(
        self, 
        baseurl: str,
        header: dict[str, str] | None = None
    ):
        self.baseurl = baseurl
        self.header = header or {}

