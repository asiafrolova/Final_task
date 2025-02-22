

GetListExpressions = function(container){
    let xhr = new XMLHttpRequest();
    xhr.open('GET','http://localhost:8080/api/v1/expressions');
    xhr.send();
    xhr.onreadystatechange = function(){
    if (xhr.readyState===4){
        if (xhr.status==200){
            
            let res=JSON.parse(xhr.responseText).expressions
            
            res.forEach(exp => {
                let elem=document.createElement("div")
                elem.classList=["Expression"]
                if (exp.status=="Completed"){
                    elem.innerHTML=`<p>Id: ${exp.id},Статус: ${exp.status}, Результат: ${exp.result}</p>`
                }else{
                    elem.innerHTML=`<p>Id: ${exp.id},Статус: ${exp.status}</p>`
                }
                console.log(exp)
                container.append(elem)
            });
                
                
            
            
            
        }else{
            alert(`Произошла ошибка: ${xhr.responseText}`)

        }
        
    }
}
}
GetExpressionByID = function(id){
    let xhr = new XMLHttpRequest();
    xhr.open('GET',`http://localhost:8080/api/v1/expressions/${id}`);
    xhr.send();
    xhr.onreadystatechange = function(){
    if (xhr.readyState===4){
        if (xhr.status==200){
            data = JSON.parse(xhr.responseText)
            console.log(xhr.responseText)
            if (data.expression.status=="Completed"){
                alert(`Выражение посчитано, результат: ${data.expression.result}`)
            }else{
                alert(`Статус обработки выражения: ${data.expression.status}`)
            }
            
        }else{
            alert(`Произошла ошибка: ${xhr.responseText}`)
        }
        
    }
}
}
PostExpression = function(exp){
    data={
        expression:exp
    }
    let xhr = new XMLHttpRequest();
    //xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    xhr.open('POST',`http://localhost:8080/api/v1/calculate`);
    xhr.send(JSON.stringify(data));
    xhr.onreadystatechange = function(){
    if (xhr.readyState===4){
        if (xhr.status==200){
            
            data = JSON.parse(xhr.responseText)
            alert(`Выражение успешно отправлено с id: ${data.id}`)
        }else{
            alert(`Произошла ошибка: ${xhr.responseText}`)

        }
        
    }
}
}


let PostExp = document.getElementById("btn-calculate")
let GetExpId = document.getElementById("btn-getByID")
let GetListExp = document.getElementById("btn-getList")
let Container = document.getElementById("container")
PostExp.addEventListener('click',e=>{
    Container.innerHTML=PostTemplate
    const form=document.getElementById("formExp")
    form.addEventListener('submit',e=>{
        e.preventDefault();
        const formData = new FormData(form);
        PostExpression(formData.get("expression"))

    })
    document.getElementById("back").addEventListener('click',e=>{
        Container.innerHTML=BaseTemplae
        location.reload()
    })
    
})
GetExpId.addEventListener('click',e=>{
    Container.innerHTML=GetByIdTemplate
    const form=document.getElementById("formExp")
    form.addEventListener('submit',e=>{
        e.preventDefault();
        const formData = new FormData(form);
        GetExpressionByID(formData.get("id"))

    })
    document.getElementById("back").addEventListener('click',e=>{
        Container.innerHTML=BaseTemplae
        location.reload()
    })
})
GetListExp.addEventListener('click',e=>{
    GetListExpressions(Container)
  
    Container.innerHTML=GetListTemplate
    document.getElementById("back").addEventListener('click',e=>{
        Container.innerHTML=BaseTemplae
        location.reload()
    })
})
BaseTemplae=`

        <button class="Btn" id="btn-calculate">
         Calculate expression
        </button>
        <button class="Btn" id="btn-getByID">
            
            Get expression by ID
        </button>
        <button class="Btn" id="btn-getList">
            Get list expressions

        </button>
    
`
PostTemplate=`

        <form id="formExp">
            <label for="expression">Введите выражениe</label>
            <input type="text" name="expression" placeholder="2+2">
            <input type="submit">
          </form>
        <button class="Btn" id="back">
         Back
        </button>
        
   
`
GetByIdTemplate=`
          <form id="formExp">
            <label for="id">Введите id выражения</label>
            <input type="text" name="id" placeholder="id1">
            <input type="submit">
          </form>
        <button class="Btn" id="back">
         Back
        </button>
`
GetListTemplate=`
<button class="Btn" id="back">
         Back
        </button>
`

