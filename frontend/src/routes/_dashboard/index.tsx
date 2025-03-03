import { getMyInfo } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { toast } from "sonner";
export const Route = createFileRoute("/_dashboard/")({
  component: RouteComponent,
});

function RouteComponent() {
  const {
    data: myInfo,
    isPending,
    isError,
    error,
  } = useQuery({
    queryKey: ["my-info"],
    queryFn: () => getMyInfo().then((res) => res.data.data),
  });

  if (isPending) return null;

  if (isError) {
    toast.error(error.message);
    return null;
  }

  const icsURL = "https://" + import.meta.env.DOMAIN + "/api/" + myInfo.id + ".ics"

  return (
    <div className="px-4 flex flex-col gap-2 mt-8 items-center">
      <h1 className="text-2xl font-bold">主页</h1>
      <span className="text-sm text-muted-foreground">
        如果你看到这里一片空白，不用担心，目前主页没有任何内容 :)
      </span>
      <span className="text-sm text-muted-foreground">
        临时将值班日程订阅入口放在这里，这个链接可以直接下载 .ics 日历数据交换标准文件
      </span>
      <span className="text-sm text-muted-foreground">
        可以在日历软件中导入，也可以在日历软件中订阅这个链接，推荐在手机原生日历软件导入或订阅
      </span>
      <a href={icsURL}>
        {icsURL}
      </a>
    </div>
  );
}
